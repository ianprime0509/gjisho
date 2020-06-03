// Package jmdict contains types and functions for working with JMdict data.
package jmdict

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ianprime0509/gjisho/xmlutil"
)

// JMdict is the JMdict database, containing data on Japanese words and phrases.
type JMdict struct {
	db          *sql.DB
	lookupQuery *sql.Stmt
	fetchQuery  *sql.Stmt
}

// New returns a new JMdict using the given database.
func New(db *sql.DB) (*JMdict, error) {
	lookupQuery, err := db.Prepare("SELECT heading, primary_reading, gloss_summary, id FROM Lookup WHERE Lookup MATCH ?")
	if err != nil {
		return nil, fmt.Errorf("could not prepare JMdict lookup query: %v", err)
	}
	fetchQuery, err := db.Prepare("SELECT data FROM Entry WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("could not prepare JMdict fetch query: %v", err)
	}

	return &JMdict{db, lookupQuery, fetchQuery}, nil
}

// ConvertInto converts the JMdict data from XML into the given database.
func ConvertInto(xmlPath string, db *sql.DB) error {
	log.Print("Converting JMdict to database")
	entities, err := xmlutil.ParseEntities(xmlPath)
	if err != nil {
		return fmt.Errorf("could not parse XML entities: %v", err)
	}

	jmdict, err := os.Open(xmlPath)
	if err != nil {
		return fmt.Errorf("could not open JMdict file: %v", err)
	}
	defer jmdict.Close()

	if err := createTables(db); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	insertEntry, err := tx.Prepare("INSERT INTO Entry VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare Entry insert statement: %v", err)
	}
	insertLookup, err := tx.Prepare("INSERT INTO Lookup VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare Lookup insert statement: %v", err)
	}

	done := 0
	decoder := xml.NewDecoder(bufio.NewReader(jmdict))
	decoder.Entity = entities
	tok, err := decoder.Token()
	for err == nil {
		if start, ok := tok.(xml.StartElement); ok && start.Name.Local == "entry" {
			if err := convertEntry(decoder, &start, insertEntry, insertLookup); err != nil {
				return fmt.Errorf("could not process JMdict entry: %v", err)
			}
			done++

			if done%10000 == 0 {
				log.Printf("Done: %v\n", done)
			}
		}
		tok, err = decoder.Token()
	}
	if err != io.EOF {
		return fmt.Errorf("could not read from JMdict file: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	if _, err := db.Exec("ANALYZE"); err != nil {
		return fmt.Errorf("could not analyze database: %v", err)
	}
	return nil
}

// createTables creates the tables required for the JMdict SQLite database.
func createTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE Entry (
		id   INTEGER PRIMARY KEY,
		data BLOB NOT NULL        -- Entry data in JSON format
	)`)
	if err != nil {
		return fmt.Errorf("could not create JMdict entry table: %v", err)
	}

	_, err = db.Exec(`CREATE VIRTUAL TABLE Lookup USING FTS5 (
		heading,
		primary_reading,
		gloss_summary,
		id UNINDEXED
	)`)
	if err != nil {
		return fmt.Errorf("could not create JMdict lookup table: %v", err)
	}

	return nil
}

func convertEntry(decoder *xml.Decoder, start *xml.StartElement, insertEntry *sql.Stmt, insertLookup *sql.Stmt) error {
	var entry Entry
	if err := decoder.DecodeElement(&entry, start); err != nil {
		return fmt.Errorf("could not unmarshal entry XML: %v", err)
	}
	data, err := json.Marshal(&entry)
	if err != nil {
		return fmt.Errorf("could not marshal entry JSON: %v", err)
	}

	_, err = insertEntry.Exec(entry.ID, data)
	if err != nil {
		return fmt.Errorf("could not insert Entry data: %v", err)
	}

	_, err = insertLookup.Exec(entry.Heading(), entry.PrimaryReading(), entry.GlossSummary(), entry.ID)
	if err != nil {
		return fmt.Errorf("could not insert Lookup data: %v", err)
	}

	return nil
}

// Fetch returns the dictionary entry with the given ID.
func (dict *JMdict) Fetch(id int) (Entry, error) {
	var data []byte
	if err := dict.fetchQuery.QueryRow(id).Scan(&data); err != nil {
		return Entry{}, fmt.Errorf("scan error: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, fmt.Errorf("could not unmarshal data: %v", err)
	}
	return entry, nil
}

// Lookup looks up dictionary entries according to the given query.
func (dict *JMdict) Lookup(query string) ([]LookupResult, error) {
	if query == "" {
		return nil, nil
	}

	rows, err := dict.lookupQuery.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var results []LookupResult
	for rows.Next() {
		var result LookupResult
		if err := rows.Scan(&result.Heading, &result.PrimaryReading, &result.GlossSummary, &result.ID); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return results, nil
}

// LookupResult is the result of a dictionary lookup.
type LookupResult struct {
	Heading        string
	PrimaryReading string
	GlossSummary   string
	ID             int
}

// Entry is a single entry in the JMdict dictionary.
type Entry struct {
	ID            int            `xml:"ent_seq"`
	KanjiReadings []KanjiReading `xml:"k_ele"`
	KanaReadings  []KanaReading  `xml:"r_ele"`
	Senses        []Sense        `xml:"sense"`
}

// Heading returns the primary heading of the entry for presentation purposes.
// This is either the first kanji reading or, if there are no kanji readings,
// the first kana reading.
func (e Entry) Heading() string {
	if len(e.KanjiReadings) > 0 {
		return e.KanjiReadings[0].Reading
	}
	return e.KanaReadings[0].Reading
}

// PrimaryReading returns the primary reading of the entry (the first kana
// reading).
func (e Entry) PrimaryReading() string {
	return e.KanaReadings[0].Reading
}

// GlossSummary returns a summary of the glosses of the entry.
func (e Entry) GlossSummary() string {
	var glosses []string
	for _, sense := range e.Senses {
		for _, gloss := range sense.Glosses {
			if gloss.Language != "" {
				continue
			}
			glosses = append(glosses, gloss.Gloss)
		}
	}
	return strings.Join(glosses, "; ")
}

// KanjiReading is a reading for an entry using kanji or other non-kana
// characters.
type KanjiReading struct {
	Reading  string   `xml:"keb"`
	Info     []string `xml:"ke_inf"`
	Priority []string `xml:"ke_pri"`
}

// KanaReading is a reading for an entry using kana.
type KanaReading struct {
	Reading      string   `xml:"reb"`
	NoKanji      NoKanji  `xml:"re_nokanji"`
	Restrictions []string `xml:"re_restr"`
	Info         []string `xml:"re_inf"`
	Priority     []string `xml:"re_pri"`
}

// NoKanji is a boolean indicating whether a kana reading is not a "true"
// reading of the kanji.
type NoKanji bool

// UnmarshalXML unmarshals a NoKanji from XML. This always returns true, since
// the element will be omitted if the value is intended to be false.
func (nk *NoKanji) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if err := d.Skip(); err != nil {
		return err
	}
	*nk = true
	return nil
}

// Sense is a sense of a dictionary entry.
type Sense struct {
	KanjiRestrictions []string     `xml:"stagk"`
	KanaRestrictions  []string     `xml:"stagr"`
	CrossReferences   []string     `xml:"xref"`
	Antonyms          []string     `xml:"ant"`
	PartsOfSpeech     []string     `xml:"pos"`
	Fields            []string     `xml:"field"`
	Misc              []string     `xml:"misc"`
	LoanSources       []LoanSource `xml:"lsource"`
	Dialects          []string     `xml:"dial"`
	Glosses           []Gloss      `xml:"gloss"`
	Info              []string     `xml:"s_inf"`
}

// LoanSource is a description of the source of a loan word.
type LoanSource struct {
	Source             string             `xml:",chardata"`
	Language           Language           `xml:"lang,attr"`
	PartialDescription PartialDescription `xml:"ls_type,attr"`
	Wasei              Wasei              `xml:"ls_wasei,attr"`
}

// Language is a three-letter language code from the ISO 639-2 standard.
type Language string

// PartialDescription is a boolean indicating whether the loan source only
// partially describes the source of the associated word or phrase.
type PartialDescription bool

// UnmarshalXMLAttr unmarshals a PartialDescription from an XML attribute.
func (pd *PartialDescription) UnmarshalXMLAttr(a xml.Attr) error {
	*pd = a.Value == "part"
	return nil
}

// Wasei is a boolean indicating whether the entry is "wasei" (made from foreign
// language components but not an actual word or phrase in that language).
type Wasei bool

// UnmarshalXMLAttr unmarshals a Wasei from an XML attribute.
func (w *Wasei) UnmarshalXMLAttr(a xml.Attr) error {
	*w = a.Value == "y"
	return nil
}

// Gloss is a gloss of a word or phrase in another language.
type Gloss struct {
	Gloss    string   `xml:",chardata"`
	Language Language `xml:"lang,attr"`
	Gender   string   `xml:"g_gend,attr"`
	Type     string   `xml:"g_type,attr"`
}

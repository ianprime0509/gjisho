// Package jmdict contains types and functions for working with JMdict data.
package jmdict

import (
	"bufio"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/ianprime0509/gjisho/internal/util"
)

// JMdict is the JMdict database, containing data on Japanese words and phrases.
type JMdict struct {
	db          *sql.DB
	lookupQuery *sql.Stmt
	fetchQuery  *sql.Stmt
}

// New returns a new JMdict using the given database.
func New(db *sql.DB) (*JMdict, error) {
	lookupQuery, err := db.Prepare(`SELECT heading, primary_reading, gloss_summary, all_writings, priority, id
	FROM EntryLookup
	WHERE EntryLookup MATCH ?
	ORDER BY -bm25(EntryLookup, 10, 4, 2) + 2 * priority DESC`)
	if err != nil {
		return nil, fmt.Errorf("could not prepare JMdict lookup query: %v", err)
	}
	fetchQuery, err := db.Prepare("SELECT data FROM Entry WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("could not prepare JMdict fetch query: %v", err)
	}

	return &JMdict{db, lookupQuery, fetchQuery}, nil
}

// ConvertInto converts the JMdict data from XML into the given database. The
// given progress callback, if non-nil, is called after every 10,000th converted
// record with the total number of records converted so far.
func ConvertInto(xmlPath string, db *sql.DB, progressCB func(int)) error {
	entities, err := util.ParseEntities(xmlPath)
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
	insertLookup, err := tx.Prepare("INSERT INTO EntryLookup VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare EntryLookup insert statement: %v", err)
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

			if done%10000 == 0 && progressCB != nil {
				progressCB(done)
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

	_, err = db.Exec(`CREATE VIRTUAL TABLE EntryLookup USING FTS5 (
		heading,
		primary_reading,
		gloss_summary,
		all_writings,
		priority UNINDEXED,
		id UNINDEXED,
		prefix = '1 2 3 4 5'
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
	data, err := util.MarshalCompressed(&entry)
	if err != nil {
		return fmt.Errorf("could not marshal entry JSON: %v", err)
	}

	_, err = insertEntry.Exec(entry.ID, data)
	if err != nil {
		return fmt.Errorf("could not insert Entry data: %v", err)
	}

	_, err = insertLookup.Exec(entry.Heading(), entry.PrimaryReading(), entry.GlossSummary(),
		entry.allWritings(), entry.Priority(), entry.ID)
	if err != nil {
		return fmt.Errorf("could not insert EntryLookup data: %v", err)
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
	if err := util.UnmarshalCompressed(data, &entry); err != nil {
		return Entry{}, fmt.Errorf("could not unmarshal data: %v", err)
	}

	return entry, nil
}

// Lookup looks up dictionary entries according to the given query. The results
// are sorted such that the ones deemed most relevant to the query come first.
func (dict *JMdict) Lookup(query string) ([]LookupResult, error) {
	if query == "" {
		return nil, nil
	}

	results, err := dict.lookupRaw(convertQuery(query))
	if err != nil {
		return nil, err
	}

	return results, nil
}

// LookupByRef returns the lookup result for the entry that most closely matches
// the given reference. The reference is expected to follow the standard format
// for cross references, using ・ as a separator between parts.
func (dict *JMdict) LookupByRef(ref string) (LookupResult, error) {
	parts := strings.Split(ref, "・")
	head := parts[0]
	reading := ""
	if len(parts) > 1 {
		reading = parts[1]
	}

	results, err := dict.lookupRaw(fmt.Sprintf(`"%v"`, strings.ReplaceAll(parts[0], `"`, `""`)))
	if err != nil {
		return LookupResult{}, fmt.Errorf("could not search for reference %q: %v", ref, err)
	}

	match := LookupResult{}
	score := 0
	for _, r := range results {
		newScore := 0
		if r.Heading == head {
			newScore++
		}
		ws := strings.Fields(r.allWritings)
		for _, w := range ws {
			if w == head {
				newScore++
			}
			if w == reading {
				newScore++
			}
		}

		if newScore > score {
			match = r
			score = newScore
		}
	}
	if match.ID == 0 {
		return LookupResult{}, fmt.Errorf("no results for reference %q", ref)
	}

	return match, nil
}

// lookupRaw is the same as Lookup, but uses a raw FTS5 query rather than a
// converted one and does not sort the results in any particular order.
func (dict *JMdict) lookupRaw(query string) ([]LookupResult, error) {
	rows, err := dict.lookupQuery.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var results []LookupResult
	for rows.Next() {
		var result LookupResult
		if err := rows.Scan(&result.Heading, &result.PrimaryReading, &result.GlossSummary,
			&result.allWritings, &result.Priority, &result.ID); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return results, nil
}

// convertQuery converts the given query string into one usable by the SQLite
// FTS5 engine.
func convertQuery(query string) string {
	sb := new(strings.Builder)
	for _, token := range strings.Fields(query) {
		if sb.Len() > 0 {
			sb.WriteString(" OR ")
		}
		esc := strings.ReplaceAll(token, `"`, `""`)
		sb.WriteRune('"')
		sb.WriteString(esc)
		sb.WriteString(`"* OR "`)
		// By also including the literal string as an additional query, we give
		// greater weight to exact matches
		sb.WriteString(esc)
		sb.WriteRune('"')
	}
	return sb.String()
}

// LookupResult is the result of a dictionary lookup.
type LookupResult struct {
	Heading        string
	PrimaryReading string
	GlossSummary   string
	allWritings    string
	Priority       int
	ID             int
}

// Entry is a single entry in the JMdict dictionary.
type Entry struct {
	ID            int            `xml:"ent_seq"`
	KanjiWritings []KanjiWriting `xml:"k_ele"`
	KanaWritings  []KanaWriting  `xml:"r_ele"`
	Senses        []Sense        `xml:"sense"`
}

// Heading returns the primary heading of the entry for presentation purposes.
// This is either the first kanji writing or, if there are no kanji writings,
// the first kana writing.
func (e Entry) Heading() string {
	if len(e.KanjiWritings) > 0 {
		return e.KanjiWritings[0].Writing
	}
	return e.KanaWritings[0].Writing
}

// PrimaryReading returns the primary reading of the entry (the first kana
// writing).
func (e Entry) PrimaryReading() string {
	return e.KanaWritings[0].Writing
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

// AssociatedKanji returns all the kanji associated with the entry (e.g. because
// they are part of the entry's writing).
func (e Entry) AssociatedKanji() []string {
	// We want the set to be in order, so the value of the map is the index of the
	// element
	set := make(map[rune]int)
	idx := 0
	for _, r := range e.KanjiWritings {
		for _, c := range r.Writing {
			if unicode.Is(unicode.Han, c) {
				if _, ok := set[c]; !ok {
					set[c] = idx
					idx++
				}
			}
		}
	}

	kanji := make([]string, len(set))
	for c, idx := range set {
		kanji[idx] = string(c)
	}
	return kanji
}

// Priority returns a measure of the "priority" of the entry, where higher
// numbers indicate a more common word (with 0 being the lowest).
func (e Entry) Priority() int {
	pri := 0
	// The current implementation is rather stupid and doesn't take into account
	// the differences between the various priority lists
	for _, k := range e.KanjiWritings {
		pri += len(k.Priority)
	}
	for _, k := range e.KanaWritings {
		pri += len(k.Priority)
	}
	return pri
}

// allWritings returns a space-separated string containing all the writings of
// the entry for lookup purposes.
func (e Entry) allWritings() string {
	readings := make([]string, 0, len(e.KanjiWritings)+len(e.KanaWritings))
	for _, w := range e.KanjiWritings {
		readings = append(readings, w.Writing)
	}
	for _, w := range e.KanaWritings {
		readings = append(readings, w.Writing)
	}
	return strings.Join(readings, " ")
}

// KanjiWriting is a writing for an entry using kanji or other non-kana
// characters.
type KanjiWriting struct {
	Writing  string   `xml:"keb"`
	Info     []string `xml:"ke_inf"`
	Priority []string `xml:"ke_pri"`
}

// KanaWriting is a writing for an entry using kana.
type KanaWriting struct {
	Writing      string   `xml:"reb"`
	NoKanji      NoKanji  `xml:"re_nokanji"`
	Restrictions []string `xml:"re_restr"`
	Info         []string `xml:"re_inf"`
	Priority     []string `xml:"re_pri"`
}

// NoKanji is a boolean indicating whether a kana writing is not a "true"
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

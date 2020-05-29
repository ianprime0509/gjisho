package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// JMdict is the JMdict database, containing data on Japanese words and phrases.
type JMdict struct {
	db          *sql.DB
	lookupQuery *sql.Stmt
	fetchQuery  *sql.Stmt
}

// OpenJMdict opens the JMdict database at the given path.
func OpenJMdict(path string) (*JMdict, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("could not open JMdict database: %v", err)
	}

	lookupQuery, err := db.Prepare("SELECT key, type, heading, gloss_summary, id FROM Lookup WHERE key MATCH ?")
	if err != nil {
		return nil, fmt.Errorf("could not prepare JMdict lookup query: %v", err)
	}
	fetchQuery, err := db.Prepare("SELECT data FROM Entry WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("could not prepare JMdict fetch query: %v", err)
	}

	return &JMdict{db, lookupQuery, fetchQuery}, nil
}

// Close closes the database.
func (dict *JMdict) Close() error {
	return dict.db.Close()
}

// Fetch returns the dictionary entry with the given ID.
func (dict *JMdict) Fetch(id int) (DictEntry, error) {
	var data []byte
	if err := dict.fetchQuery.QueryRow(id).Scan(&data); err != nil {
		return DictEntry{}, fmt.Errorf("scan error: %v", err)
	}

	var entry DictEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return DictEntry{}, fmt.Errorf("could not unmarshal data: %v", err)
	}
	return entry, nil
}

// Lookup looks up dictionary entries according to the given query.
func (dict *JMdict) Lookup(query string) ([]LookupEntry, error) {
	if query == "" {
		return nil, nil
	}

	rows, err := dict.lookupQuery.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var entries []LookupEntry
	for rows.Next() {
		var entry LookupEntry
		if err := rows.Scan(&entry.Key, &entry.Type, &entry.Heading, &entry.GlossSummary, &entry.ID); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return entries, nil
}

// LookupEntry is the result of a dictionary lookup.
type LookupEntry struct {
	Key          string
	Type         string
	Heading      string
	GlossSummary string
	ID           int
}

// DictEntry is a single entry in the JMdict dictionary.
type DictEntry struct {
	ID            int            `xml:"ent_seq"`
	KanjiReadings []KanjiReading `xml:"k_ele"`
	KanaReadings  []KanaReading  `xml:"r_ele"`
	Senses        []Sense        `xml:"sense"`
}

// Heading returns the primary heading of the entry for presentation purposes.
// This is either the first kanji reading or, if there are no kanji readings,
// the first kana reading.
func (e DictEntry) Heading() string {
	if len(e.KanjiReadings) > 0 {
		return e.KanjiReadings[0].Reading
	}
	return e.KanaReadings[0].Reading
}

// GlossSummary returns a summary of the glosses of the entry.
func (e DictEntry) GlossSummary() string {
	// For conciseness, we take only the first five glosses
	var glosses []string
outer:
	for _, sense := range e.Senses {
		for _, gloss := range sense.Glosses {
			if len(glosses) >= 5 {
				break outer
			}
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

// Wasei is a boolean indicating whether the entry is "wasei" (made from foreign
// language components but not an actual word or phrase in that language).
type Wasei bool

// Gloss is a gloss of a word or phrase in another language.
type Gloss struct {
	Gloss    string   `xml:",chardata"`
	Language Language `xml:"lang,attr"`
	Gender   string   `xml:"g_gend,attr"`
	Type     string   `xml:"g_type,attr"`
}

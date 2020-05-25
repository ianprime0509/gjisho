package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	lookupQuery, err := db.Prepare("SELECT key, type, id FROM Lookup WHERE key LIKE ?")
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

// Fetch returns the dictionary entries with the given IDs.
func (dict *JMdict) Fetch(ids []int) ([]DictEntry, error) {
	var entries []DictEntry
	for _, id := range ids {
		row := dict.fetchQuery.QueryRow(id)

		var data []byte
		if err := row.Scan(&data); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		var entry DictEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			return nil, fmt.Errorf("could not unmarshal data: %v", err)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Lookup returns the IDs of all dictionary entries corresponding to the given
// query.
func (dict *JMdict) Lookup(query string) ([]int, error) {
	rows, err := dict.lookupQuery.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var _key, _typ string
		var id int
		if err := rows.Scan(&_key, &_typ, &id); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		ids = append(ids, id)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return ids, nil
}

// DictEntry is a single entry in the JMdict dictionary.
type DictEntry struct {
	ID            int            `xml:"ent_seq"`
	KanjiReadings []KanjiReading `xml:"k_ele"`
	KanaReadings  []KanaReading  `xml:"r_ele"`
	Senses        []Sense        `xml:"sense"`
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
	Language           Language           `xml:"xml:lang,attr"`
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
	Gloss  string `xml:",chardata"`
	Gender string `xml:"g_gend,attr"`
	Type   string `xml:"g_type,attr"`
}

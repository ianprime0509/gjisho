package dbconv

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// DictEntry is a single entry in the JMdict dictionary.
type DictEntry struct {
	Sequence      int            `xml:"ent_seq"`
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
	Language           Language           `xml:"xml:lang,attr"`
	PartialDescription PartialDescription `xml:"ls_type,attr"`
	Wasei              Wasei              `xml:"ls_wasei,attr"`
}

// Language is a three-letter language code from the ISO 639-2 standard.
type Language string

// UnmarshalXMLAttr unmarshals a Language from an XML attribute.
func (lang *Language) UnmarshalXMLAttr(a xml.Attr) error {
	if a.Value == "" {
		*lang = "eng"
	} else if len(a.Value) != 3 {
		return xml.UnmarshalError(fmt.Sprintf("invalid language code: %v", a.Value))
	} else {
		*lang = Language(a.Value)
	}
	return nil
}

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
	Gloss  string `xml:",chardata"`
	Gender string `xml:"g_gend,attr"`
	Type   string `xml:"g_type,attr"`
}

// ConvertJMdict converts the JMdict data from XML to a SQLite database.
func ConvertJMdict(xmlPath string, dbPath string) error {
	log.Println("Converting JMdict to database")
	entities, err := parseEntities(xmlPath)
	if err != nil {
		return fmt.Errorf("could not parse XML entities: %v", err)
	}

	jmdict, err := os.Open(xmlPath)
	if err != nil {
		return fmt.Errorf("could not open JMdict file: %v", err)
	}
	defer jmdict.Close()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("could not open database: %v", err)
	}
	defer db.Close()

	if err := createJMdictTables(db); err != nil {
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
	insertLookup, err := tx.Prepare("INSERT INTO Lookup VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare Lookup insert statement: %v", err)
	}

	done := 0
	decoder := xml.NewDecoder(bufio.NewReader(jmdict))
	decoder.Entity = entities
	tok, err := decoder.Token()
	for err == nil {
		if start, ok := tok.(xml.StartElement); ok && start.Name.Local == "entry" {
			if err := convertDictEntry(decoder, &start, insertEntry, insertLookup); err != nil {
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

var entityRegexp = regexp.MustCompile(`<!ENTITY +([^" ]+) *"([^"]*)">`)
var endDoctypeRegexp = regexp.MustCompile(`]>`)

// parseEntities parses the entity definitions from the doctype of the given XML
// file. It is not very smart and does only enough to be useful.
func parseEntities(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()
	input := bufio.NewReader(file)

	entities := make(map[string]string)
	line, err := input.ReadString('\n')
	for err == nil {
		if endDoctypeRegexp.MatchString(line) {
			return entities, nil
		}
		if matches := entityRegexp.FindStringSubmatch(line); matches != nil {
			entities[matches[1]] = matches[2]
		}

		line, err = input.ReadString('\n')
	}

	if err != io.EOF {
		return nil, fmt.Errorf("could not read from file: %v", err)
	}
	return entities, nil
}

// createJMdictTables creates the tables required for the JMdict SQLite
// database.
func createJMdictTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE Entry (
		id   INTEGER PRIMARY KEY,
		data BLOB NOT NULL        -- Entry data in JSON format
	)`)
	if err != nil {
		return fmt.Errorf("could not create JMdict entry table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE Lookup (
		key TEXT NOT NULL COLLATE NOCASE,
		id  INTEGER NOT NULL REFERENCES Entry(id)
	)`)
	if err != nil {
		return fmt.Errorf("could not create JMdict lookup table: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX Lookup_key ON Lookup(key)`)
	if err != nil {
		return fmt.Errorf("could not create JMdict lookup index: %v", err)
	}

	return nil
}

func convertDictEntry(decoder *xml.Decoder, start *xml.StartElement, insertEntry *sql.Stmt, insertLookup *sql.Stmt) error {
	var entry DictEntry
	if err := decoder.DecodeElement(&entry, start); err != nil {
		return fmt.Errorf("could not unmarshal entry XML: %v", err)
	}
	data, err := json.Marshal(&entry)
	if err != nil {
		return fmt.Errorf("could not marshal entry JSON: %v", err)
	}

	_, err = insertEntry.Exec(entry.Sequence, data)
	if err != nil {
		return fmt.Errorf("could not insert Entry data: %v", err)
	}
	for _, kanji := range entry.KanjiReadings {
		_, err = insertLookup.Exec(kanji.Reading, entry.Sequence)
		if err != nil {
			return fmt.Errorf("could not insert Lookup data for kanji: %v", err)
		}
	}
	for _, kana := range entry.KanaReadings {
		_, err = insertLookup.Exec(kana.Reading, entry.Sequence)
		if err != nil {
			return fmt.Errorf("could not insert Lookup data for kana: %v", err)
		}
	}
	for _, sense := range entry.Senses {
		for _, gloss := range sense.Glosses {
			_, err := insertLookup.Exec(gloss.Gloss, entry.Sequence)
			if err != nil {
				return fmt.Errorf("could not insert Lookup data for gloss: %v", err)
			}
		}
	}

	return nil
}

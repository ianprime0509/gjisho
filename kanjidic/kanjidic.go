// Package kanjidic contains types and functions for working with Kanjidic2
// data.
package kanjidic

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ianprime0509/gjisho/xmlutil"
)

// Kanjidic is the Kanjidic2 database, containing data on kanji.
type Kanjidic struct {
	db         *sql.DB
	fetchQuery *sql.Stmt
}

// New returns a new Kanjidic using the given database.
func New(db *sql.DB) (*Kanjidic, error) {
	fetchQuery, err := db.Prepare("SELECT data FROM Kanji WHERE character = ?")
	if err != nil {
		return nil, fmt.Errorf("could not prepare Kanjidic fetch query: %v", err)
	}
	return &Kanjidic{db, fetchQuery}, nil
}

// ConvertInto converts the Kanjidic2 data from XML into the given database.
func ConvertInto(xmlPath string, db *sql.DB) error {
	log.Print("Converting Kanjidic2 to database")
	entities, err := xmlutil.ParseEntities(xmlPath)
	if err != nil {
		return fmt.Errorf("could not parse XML entities: %v", err)
	}

	kanjidic, err := os.Open(xmlPath)
	if err != nil {
		return fmt.Errorf("could not open Kanjidic file: %v", err)
	}
	defer kanjidic.Close()

	if err := createTables(db); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	insert, err := tx.Prepare("INSERT INTO Kanji VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare Kanji insert statement: %v", err)
	}

	done := 0
	decoder := xml.NewDecoder(bufio.NewReader(kanjidic))
	decoder.Entity = entities
	tok, err := decoder.Token()
	for err == nil {
		if start, ok := tok.(xml.StartElement); ok && start.Name.Local == "character" {
			if err := convertEntry(decoder, &start, insert); err != nil {
				return fmt.Errorf("could not process Kanjidic entry: %v", err)
			}
			done++

			if done%1000 == 0 {
				log.Printf("Done: %v\n", done)
			}
		}
		tok, err = decoder.Token()
	}
	if err != io.EOF {
		return fmt.Errorf("could not read from Kanjidic file: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	if _, err := db.Exec("ANALYZE"); err != nil {
		return fmt.Errorf("could not analyze database: %v", err)
	}
	return nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE Kanji (
		character TEXT PRIMARY KEY,
		data BLOB NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("could not create kanji table: %v", err)
	}

	return nil
}

func convertEntry(decoder *xml.Decoder, start *xml.StartElement, insert *sql.Stmt) error {
	var entry Character
	if err := decoder.DecodeElement(&entry, start); err != nil {
		return fmt.Errorf("could not unmarshal entry XML: %v", err)
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("could not marshal entry JSON: %v", err)
	}

	_, err = insert.Exec(entry.Literal, data)
	if err != nil {
		return fmt.Errorf("could not insert Kanji data: %v", err)
	}

	return nil
}

// Fetch returns the character data for the given kanji.
func (dict *Kanjidic) Fetch(kanji string) (Character, error) {
	var data []byte
	if err := dict.fetchQuery.QueryRow(kanji).Scan(&data); err != nil {
		return Character{}, fmt.Errorf("scan error: %v", err)
	}

	var char Character
	if err := json.Unmarshal(data, &char); err != nil {
		return Character{}, fmt.Errorf("could not unmarshal data: %v", err)
	}

	return char, nil
}

// Character is a single kanji character.
type Character struct {
	Literal              string                `xml:"literal"`
	Codepoints           []Codepoint           `xml:"codepoint>cp_value"`
	Radicals             []Radical             `xml:"radical>rad_value"`
	Misc                 Misc                  `xml:"misc"`
	DictRefs             []DictRef             `xml:"dic_number>dic_ref"`
	QueryCodes           []QueryCode           `xml:"query_code>q_code"`
	ReadingMeaningGroups []ReadingMeaningGroup `xml:"reading_meaning>rmgroup"`
	Nanori               []string              `xml:"reading_meaning>nanori"`
}

// Readings returns all the character's readings of the given type.
func (c Character) Readings(ty ReadingType) []string {
	var rs []string
	for _, group := range c.ReadingMeaningGroups {
		for _, r := range group.Readings {
			if r.Type == ty {
				rs = append(rs, r.Reading)
			}
		}
	}
	return rs
}

// Meanings returns all the character's English meanings.
func (c Character) Meanings() []string {
	var ms []string
	for _, group := range c.ReadingMeaningGroups {
		for _, m := range group.Meanings {
			if m.Language == "" {
				ms = append(ms, m.Meaning)
			}
		}
	}
	return ms
}

// Codepoint is a description of a character's codepoint in some standard.
type Codepoint struct {
	Value string `xml:",chardata"`
	Type  string `xml:"cp_type,attr"`
}

// Radical is a description of a character radical.
type Radical struct {
	Value string `xml:",chardata"`
	Type  string `xml:"rad_type,attr"`
}

// Misc contains miscellaneous information about a character.
type Misc struct {
	Grade        Grade     `xml:"grade"`
	StrokeCounts []int     `xml:"stroke_count"`
	Variants     []Variant `xml:"variant"`
	Frequency    int       `xml:"freq"`
	RadicalName  []string  `xml:"rad_name"`
	JLPTLevel    int       `xml:"jlpt"`
}

// Grade is a kanji grade level.
type Grade int

const (
	// None indicates a kanji that is not included in any official list.
	None = 0
	// JuniorHigh indicates a jouyou kanji taught in junior high.
	JuniorHigh = 8
	// Jinmeiyou indicates a jinmeiyou kanji.
	Jinmeiyou = 9
	// JouyouVariantJinmeiyou indicates a jinmeiyou kanji which is a variant of a
	// jouyou kanji.
	JouyouVariantJinmeiyou = 10
)

func (g Grade) String() string {
	switch g {
	case None:
		return "not in standard list"
	case JuniorHigh:
		return "taught in junior high"
	case Jinmeiyou:
		return "jinmeiyou"
	case JouyouVariantJinmeiyou:
		return "jinmeiyou (variant of jouyou)"
	default:
		return fmt.Sprintf("taught in grade %v", int(g))
	}
}

// Variant is a description of a character variant.
type Variant struct {
	Value string `xml:",chardata"`
	Type  string `xml:"var_type,attr"`
}

// DictRef is a reference to a kanji in a dictionary.
type DictRef struct {
	Index  string `xml:",chardata"`
	Type   string `xml:"dr_type,attr"`
	Volume int    `xml:"m_vol,attr"`
	Page   int    `xml:"m_page,attr"`
}

// QueryCode is a code that can be used to find a character in some system.
type QueryCode struct {
	Code         string `xml:",chardata"`
	Type         string `xml:"qc_type,attr"`
	SKIPMisclass string `xml:"skip_misclass,attr"`
}

// ReadingMeaningGroup is a group of related readings and meanings of a
// character.
type ReadingMeaningGroup struct {
	Readings []Reading `xml:"reading"`
	Meanings []Meaning `xml:"meaning"`
}

// Reading is a reading of a character.
type Reading struct {
	Reading string      `xml:",chardata"`
	Type    ReadingType `xml:"r_type,attr"`
	OnType  string      `xml:"on_type,attr"`
	Jouyou  Jouyou      `xml:"r_status,attr"`
}

// ReadingType is the type of a reading.
type ReadingType string

const (
	// Pinyin is a Chinese reading in Pinyin.
	Pinyin ReadingType = "pinyin"
	// KoreanRomanized is a romanized Korean reading.
	KoreanRomanized = "korean_r"
	// KoreanHangul is a Korean reading in Hangul.
	KoreanHangul = "korean_h"
	// Vietnamese is a Vietnamese reading.
	Vietnamese = "vietnam"
	// On is a Japanese on reading.
	On = "ja_on"
	// Kun is a Japanese kun reading.
	Kun = "ja_kun"
)

// Jouyou is a boolean value indicating whether a reading is approved for a
// Jouyou kanji.
type Jouyou bool

// UnmarshalXMLAttr unmarshals a Jouyou from an XML attribute value.
func (j *Jouyou) UnmarshalXMLAttr(attr xml.Attr) {
	*j = attr.Value == "jy"
}

// Meaning is a meaning of a character.
type Meaning struct {
	Meaning  string   `xml:",chardata"`
	Language Language `xml:"m_lang,attr"`
}

// Language is a two-letter language code from the ISO 639-1 standard.
type Language string

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
		data      BLOB NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("could not create kanji table: %v", err)
	}

	return nil
}

func convertEntry(decoder *xml.Decoder, start *xml.StartElement, insert *sql.Stmt) error {
	var entry Character
	if err := decoder.DecodeElement(&entry, start); err != nil {
		return fmt.Errorf("could not unmarshal kanji XML: %v", err)
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("could not marshal kanji JSON: %v", err)
	}

	_, err = insert.Exec(entry.Literal, data)
	if err != nil {
		return fmt.Errorf("could not insert kanji data: %v", err)
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
	Value string        `xml:",chardata"`
	Type  CodepointType `xml:"cp_type,attr"`
}

// CodepointType is a type of a codepoint.
type CodepointType string

const (
	// JIS208Codepoint is a codepoint in the JIS X 0208-1997 encoding standard.
	JIS208Codepoint CodepointType = "jis208"
	// JIS212Codepoint is a codepoint in the JIS X 0212-1990 encoding standard.
	JIS212Codepoint = "jis212"
	// JIS213Codepoint is a codepoint in the JIS X 0213-2000 encoding standard.
	JIS213Codepoint = "jis213"
	// UCSCodepoint is a codepoint in the Unicode encoding standard.
	UCSCodepoint = "ucs"
)

func (t CodepointType) String() string {
	switch t {
	case JIS208Codepoint:
		return "JIS208"
	case JIS212Codepoint:
		return "JIS212"
	case JIS213Codepoint:
		return "JIS213"
	case UCSCodepoint:
		return "UCS"
	default:
		return string(t)
	}
}

// Radical is a description of a character radical.
type Radical struct {
	Value string      `xml:",chardata"`
	Type  RadicalType `xml:"rad_type,attr"`
}

// RadicalType is a type of a radical.
type RadicalType string

const (
	// ClassicalRadical is a classical radical, as recorded in the KangXi Zidian.
	ClassicalRadical RadicalType = "classical"
	// ClassicNelsonRadical is a radical as used in the Nelson "Modern Japanese-English
	// Character Dictionary".
	ClassicNelsonRadical = "nelson_c"
)

func (t RadicalType) String() string {
	switch t {
	case ClassicalRadical:
		return "classical"
	case ClassicNelsonRadical:
		return "Classic Nelson"
	default:
		return string(t)
	}
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
	Value string      `xml:",chardata"`
	Type  VariantType `xml:"var_type,attr"`
}

// VariantType is a type of a character variant.
type VariantType string

const (
	// JIS208Variant is a variant coding in the JIS X 0208 standard.
	JIS208Variant VariantType = "jis208"
	// JIS212Variant is a variant coding in the JIS X 0212 standard.
	JIS212Variant = "jis212"
	// JIS213Variant is a variant coding in the JIS X 0213 standard.
	JIS213Variant = "jis213"
	// DeRooVariant is a variant De Roo number.
	DeRooVariant = "deroo"
	// NJECDVariant is a variant index number in the NJECD dictionary by Halpern.
	NJECDVariant = "njecd"
	// SHVariant is a variant descriptor in "The Kanji Dictionary" by Spahn and
	// Hadamitzky.
	SHVariant = "s_h"
	// ClassicNelsonVariant is a variant code in the "Classic Nelson" dictionary.
	ClassicNelsonVariant = "nelson_c"
	// ONeillVariant is a variant code in "Japanese Names" by O'Neill.
	ONeillVariant = "oneill"
	// UCSVariant is a variant UCS/Unicode codepoint.
	UCSVariant = "ucs"
)

// DictRef is a reference to a kanji in a dictionary.
type DictRef struct {
	Index  string      `xml:",chardata"`
	Type   DictRefType `xml:"dr_type,attr"`
	Volume int         `xml:"m_vol,attr"`
	Page   int         `xml:"m_page,attr"`
}

// DictRefType is a type of a dictionary reference.
type DictRefType string

const (
	// ClassicNelsonDict is the "Modern Reader's Japanese-English Character
	// Dictionary", edited by Andrew Nelson.
	ClassicNelsonDict DictRefType = "nelson_c"
	// NewNelsonDict is "The New Nelson Japanese-English Character Dictionary",
	// edited by John Haig.
	NewNelsonDict = "nelson_n"
	// HalpernNJECDDict is the "New Japanese-English Character Dictionary", edited
	// by Jack Halpern.
	HalpernNJECDDict = "halpern_njecd"
	// HalpernKKDDict is the "Kodansha Kanji Dictionary", (2nd Ed. of the NJECD)
	// edited by Jack Halpern.
	HalpernKKDDict = "halpern_kkd"
	// HalpernKKLDDict is the "Kanji Learners Dictionary" (Kodansha) edited by
	// Jack Halpern.
	HalpernKKLDDict = "halpern_kkld"
	// HalpernKKLD2Dict is the "Kanji Learners Dictionary" (Kodansha), 2nd edition
	// (2013) edited by Jack Halpern.
	HalpernKKLD2Dict = "halpern_kkld_2ed"
	// HeisigDict is "Remembering The Kanji" by James Heisig.
	HeisigDict = "heisig"
	// Heisig6Dict is "Remembering The Kanji, Sixth Ed." by James Heisig.
	Heisig6Dict = "heisig6"
	// GakkenDict is "A New Dictionary of Kanji Usage" (Gakken).
	GakkenDict = "gakken"
	// ONeillNamesDict is "Japanese Names", by P.G. O'Neill.
	ONeillNamesDict = "oneill_names"
	// ONeillKKDict is "Essential Kanji" by P.G. O'Neill.
	ONeillKKDict = "oneill_kk"
	// MorohashiDict is the "Daikanwajiten" compiled by Morohashi.
	MorohashiDict = "moro"
	// HenshallDict is "A Guide To Remembering Japanese Characters" by Kenneth G.
	// Henshall.
	HenshallDict = "henshall"
	// SHKKDict is "Kanji and Kana" by Spahn and Hadamitzky.
	SHKKDict = "sh_kk"
	// SHKK2Dict is "Kanji and Kana" by Spahn and Hadamitzky (2011 edition).
	SHKK2Dict = "sh_kk2"
	// SakadeDict is "A Guide To Reading and Writing Japanese" edited by Florence
	// Sakade.
	SakadeDict = "sakade"
	// JFCardsDict is Japanese Kanji Flashcards, by Max Hodges and Tomoko Okazaki.
	// (Series 1).
	JFCardsDict = "jf_cards"
	// Henshall3Dict is "A Guide To Reading and Writing Japanese" 3rd edition,
	// edited by Henshall, Seeley and De Groot.
	Henshall3Dict = "henshall3"
	// TuttleCardsDict is Tuttle Kanji Cards, compiled by Alexander Kask.
	TuttleCardsDict = "tutt_cards"
	// CrowleyDict is "The Kanji Way to Japanese Language Power" by Dale Crowley.
	CrowleyDict = "crowley"
	// KanjiInContextDict is "Kanji in Context" by Nishiguchi and Kono.
	KanjiInContextDict = "kanji_in_context"
	// BusyPeopleDict is "Japanese For Busy People" vols I-III, published by the
	// AJLT.
	BusyPeopleDict = "busy_people"
	// KodanshaCompactDict is the "Kodansha Compact Kanji Guide".
	KodanshaCompactDict = "kodansha_compact"
	// ManietteDict is Yves Maniette's "Les Kanjis dans la tete" French adaptation
	// of Heisig.
	ManietteDict = "maniette"
)

func (t DictRefType) String() string {
	switch t {
	case ClassicNelsonDict:
		return "Modern Reader's Japanese-English Character Dictionary"
	case NewNelsonDict:
		return "The New Nelson Japanese-English Character Dictionary"
	case HalpernNJECDDict:
		return "New Japanese-English Character Dictionary"
	case HalpernKKDDict:
		return "Kodansha Kanji Dictionary"
	case HalpernKKLDDict:
		return "Kanji Learners Dictionary"
	case HalpernKKLD2Dict:
		return "Kanji Learners Dictionary, 2nd edition"
	case HeisigDict:
		return "Remembering The Kanji"
	case Heisig6Dict:
		return "Remembering The Kanji, Sixth Ed."
	case GakkenDict:
		return "A New Dictionary of Kanji Usage"
	case ONeillNamesDict:
		return "Japanese Names"
	case ONeillKKDict:
		return "Essential Kanji"
	case MorohashiDict:
		return "Daikanwajiten"
	case HenshallDict:
		return "A Guide To Remembering Japanese Characters"
	case SHKKDict:
		return "Kanji and Kana"
	case SHKK2Dict:
		return "Kanji and Kana, 2011 edition"
	case SakadeDict:
		return "A Guide To Reading and Writing Japanese"
	case JFCardsDict:
		return "Japanese Kanji Flashcards"
	case Henshall3Dict:
		return "A Guide To Reading and Writing Japanese, 3rd edition"
	case TuttleCardsDict:
		return "Tuttle Kanji Cards"
	case CrowleyDict:
		return "The Kanji Way to Japanese Language Power"
	case KanjiInContextDict:
		return "Kanji in Context"
	case BusyPeopleDict:
		return "Japanese For Busy People"
	case KodanshaCompactDict:
		return "Kodansha Compact Kanji Guide"
	case ManietteDict:
		return "Les Kanjis dans la tete"
	default:
		return string(t)
	}
}

// QueryCode is a code that can be used to find a character in some system.
type QueryCode struct {
	Code         string        `xml:",chardata"`
	Type         QueryCodeType `xml:"qc_type,attr"`
	SKIPMisclass string        `xml:"skip_misclass,attr"`
}

// QueryCodeType is a type of a query code.
type QueryCodeType string

const (
	// SKIPCode is a code from Halpern's SKIP (System of Kanji Indexing by
	// Patterns) system.
	SKIPCode QueryCodeType = "skip"
	// SHDescCode is a code from The Kanji Dictionary by Spahn and Hadamitzky.
	SHDescCode = "sh_desc"
	// FourCornerCode is a code from the "Four Corner" system developed by Wang
	// Chen.
	FourCornerCode = "four_corner"
	// DeRooCode is a code from Joseph De Roo's book "2001 Kanji".
	DeRooCode = "deroo"
)

func (t QueryCodeType) String() string {
	switch t {
	case SKIPCode:
		return "SKIP"
	case SHDescCode:
		return "The Kanji Dictionary (Spahn and Hadamitzky)"
	case FourCornerCode:
		return "Four Corner"
	case DeRooCode:
		return "De Roo"
	default:
		return string(t)
	}
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
	OnType  OnType      `xml:"on_type,attr"`
	Jouyou  Jouyou      `xml:"r_status,attr"`
}

// ReadingType is the type of a reading.
type ReadingType string

const (
	// PinyinReading is a Chinese reading in pinyin.
	PinyinReading ReadingType = "pinyin"
	// KoreanRomanizedReading is a romanized Korean reading.
	KoreanRomanizedReading = "korean_r"
	// KoreanHangulReading is a Korean reading in Hangul.
	KoreanHangulReading = "korean_h"
	// VietnameseReading is a VietnameseReading reading.
	VietnameseReading = "vietnam"
	// OnReading is a Japanese on reading.
	OnReading = "ja_on"
	// KunReading is a Japanese kun reading.
	KunReading = "ja_kun"
)

func (t ReadingType) String() string {
	switch t {
	case PinyinReading:
		return "pinyin"
	case KoreanRomanizedReading:
		return "Korean (romanized)"
	case KoreanHangulReading:
		return "Korean (Hangul)"
	case VietnameseReading:
		return "Vietnamese"
	case OnReading:
		return "on"
	case KunReading:
		return "kun"
	default:
		return string(t)
	}
}

// OnType is the type of an on reading.
type OnType string

const (
	// KanOn is a kan on reading.
	KanOn OnType = "kan"
	// GoOn is a go on reading.
	GoOn = "go"
	// TouOn is a tou on reading.
	TouOn = "tou"
	// KanyouOn is a kan'you on reading.
	KanyouOn = "kan'you"
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

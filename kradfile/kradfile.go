// Package kradfile contains types and functions for working with data from the
// KRADFILE kanji decomposition file.
package kradfile

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/ianprime0509/gjisho/kanjidic"
)

// KRADFILE is the KRADFILE database, containing associations between kanji and
// components.
type KRADFILE struct {
	db           *sql.DB
	fetchByKanji *sql.Stmt
}

// New returns a new KRADFILE using the given database.
func New(db *sql.DB) (*KRADFILE, error) {
	fetchByKanji, err := db.Prepare("SELECT radicals FROM Radical WHERE kanji = ?")
	if err != nil {
		return nil, fmt.Errorf("could not prepare query to fetch radicals by kanji: %v", err)
	}

	return &KRADFILE{db, fetchByKanji}, nil
}

// ConvertInto converts the KRADFILE data from plain text into the given
// database.
func ConvertInto(txtPath string, db *sql.DB) error {
	kanjiDict, err := kanjidic.New(db)
	if err != nil {
		return fmt.Errorf("could not open KANJIDIC database: %v", err)
	}

	kradfile, err := os.Open(txtPath)
	if err != nil {
		return fmt.Errorf("could not open KRADFILE file: %v", err)
	}
	defer kradfile.Close()

	if err := createTables(db); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	insertAssoc, err := tx.Prepare("INSERT INTO Radical VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare Radical insert statement: %v", err)
	}

	s := bufio.NewScanner(kradfile)
	for s.Scan() {
		if err := convertAssociation(s.Text(), kanjiDict, insertAssoc); err != nil {
			return fmt.Errorf("could not convert radical association %q: %v", s.Text(), err)
		}
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("could not read from KRADFILE file: %v", err)
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
	_, err := db.Exec(`CREATE TABLE Radical (
		kanji    TEXT NOT NULL,
		radicals TEXT NOT NULL,
		strokes  INTEGER NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("could not create Radical table: %v", err)
	}

	_, err = db.Exec("CREATE INDEX Radical_kanji ON Radical(kanji)")
	if err != nil {
		return fmt.Errorf("could not create index on Radical(kanji): %v", err)
	}

	_, err = db.Exec("CREATE INDEX Radical_radicals ON Radical(radicals)")
	if err != nil {
		return fmt.Errorf("could not create index on Radical(radicals): %v", err)
	}

	return nil
}

// trueRadicals is an association between the radicals provided in the KRADFILE
// data and the actual radical represented. Due to encoding limitations,
// KRADFILE has had to use some other characters in place of certain radicals,
// but we want the actual Unicode representations of the radicals.
var trueRadicals = map[string]string{
	"化": "⺅",
	"个": "𠆢",
	"并": "丷",
	"刈": "⺉",
	"込": "⻌",
	"尚": "⺌",
	"忙": "⺖",
	"扎": "扌",
	"汁": "⺡",
	"犯": "⺨",
	"艾": "⺾",
	"邦": "⻏",
	"阡": "阝",
	"老": "⺹",
	"杰": "⺣",
	"礼": "礻",
	"疔": "⽧",
	"禹": "⽱",
	"初": "⻂",
	"買": "⺲",
	"滴": "啇",
	"乞": "𠂉",
	// Visually very similar, but not truly the radical characters
	"｜": "丨",
	"ノ": "丿",
	"マ": "龴",
	"ハ": "八",
	"ヨ": "彐",
	// Unfortunately, ユ does not seem to have an actual radical equivalent
}

func convertAssociation(assoc string, kanjiDict *kanjidic.KANJIDIC, insert *sql.Stmt) error {
	if strings.HasPrefix(assoc, "#") {
		// Comment line
		return nil
	}

	parts := strings.SplitN(assoc, " : ", 2)
	kanji := parts[0]
	kanjiDetails, err := kanjiDict.Fetch(kanji)
	if err != nil {
		return fmt.Errorf("could not fetch details for kanji %q: %v", kanji, err)
	}
	if len(kanjiDetails.Misc.StrokeCounts) == 0 {
		return fmt.Errorf("no stroke count information for %q", kanji)
	}
	strokes := kanjiDetails.Misc.StrokeCounts[0]

	rads := strings.Split(parts[1], " ")
	for i, rad := range strings.Split(parts[1], " ") {
		if sub, ok := trueRadicals[rad]; ok {
			rads[i] = sub
		}
	}
	if _, err := insert.Exec(kanji, strings.Join(rads, ""), strokes); err != nil {
		return fmt.Errorf("could not insert association between %q and %q: %v", kanji, rads, err)
	}

	return nil
}

// FetchByKanji returns all the radicals for the given kanji.
func (k *KRADFILE) FetchByKanji(kanji string) ([]string, error) {
	// TODO: update
	rows, err := k.fetchByKanji.Query(kanji)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}

	var rads []string
	for rows.Next() {
		var rad string
		if err := rows.Scan(&rad); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		rads = append(rads, rad)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return rads, nil
}

// FetchByRadicals returns all the kanji containing all of the given radicals as
// well as all the radicals used by the returned kanji. As a special case, if no
// radicals are provided, this function returns no kanji and all the radicals
// defined in KRADFILE.
func (k *KRADFILE) FetchByRadicals(rads []string) (kanji []Kanji, krads []string, err error) {
	if len(rads) == 0 {
		krads = make([]string, 0, len(RadicalStrokes))
		for rad := range RadicalStrokes {
			krads = append(krads, rad)
		}
		return nil, krads, nil
	}

	where := new(strings.Builder)
	params := make([]interface{}, 0, len(rads))
	for _, rad := range rads {
		params = append(params, "%"+rad+"%")
		if where.Len() > 0 {
			where.WriteString(" AND ")
		}
		where.WriteString("radicals LIKE ?")
	}

	stmt, err := k.db.Prepare(fmt.Sprintf("SELECT kanji, radicals, strokes FROM Radical WHERE %v", where))
	if err != nil {
		return nil, nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, nil, fmt.Errorf("query error: %v", err)
	}

	radSet := make(map[string]struct{})
	for rows.Next() {
		var k, rads string
		var strokes int
		if err := rows.Scan(&k, &rads, &strokes); err != nil {
			return nil, nil, fmt.Errorf("scan error: %v", err)
		}

		radStrings := strings.Split(rads, "")
		kanji = append(kanji, Kanji{
			Literal:     k,
			Radicals:    radStrings,
			StrokeCount: strokes,
		})
		for _, rad := range radStrings {
			radSet[rad] = struct{}{}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("rows error: %v", err)
	}

	krads = make([]string, 0, len(radSet))
	for k := range radSet {
		krads = append(krads, k)
	}

	return kanji, krads, nil
}

// Kanji is a kanji with its associated radicals and stroke count.
type Kanji struct {
	Literal     string
	Radicals    []string
	StrokeCount int
}

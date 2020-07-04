// Package tatoeba contains types and functions for working with data from the
// Tatoeba collection of example sentences.
//
// A description of the Tatoeba file format can be found on the EDRDG website:
// https://www.edrdg.org/wiki/index.php/Sentence-Dictionary_Linking
package tatoeba

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/ianprime0509/gjisho/internal/util"
)

// Tatoeba is the Tatoeba database, containing Japanese-English example
// sentences.
type Tatoeba struct {
	db               *sql.DB
	fetchByWordQuery *sql.Stmt
}

// New returns a new Tatoeba using the given database.
func New(db *sql.DB) (*Tatoeba, error) {
	fetchByWordQuery, err := db.Prepare(`SELECT data
	FROM Example
	WHERE id IN (SELECT DISTINCT example_id FROM ExampleLookup WHERE word = ?)`)
	if err != nil {
		return nil, fmt.Errorf("could not prepare Tatoeba fetch by word query: %v", err)
	}
	return &Tatoeba{db, fetchByWordQuery}, nil
}

// ConvertInto converts the Tatoeba data from plain text into the given
// database. The given progress callback, if non-nil, is called after every
// 10,000th record with the total number of records converted so far.
func ConvertInto(txtPath string, db *sql.DB, progressCB func(int)) error {
	tatoeba, err := os.Open(txtPath)
	if err != nil {
		return fmt.Errorf("could not open Tatoeba file: %v", err)
	}
	defer tatoeba.Close()

	if err := createTables(db); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	insertExample, err := tx.Prepare("INSERT INTO Example VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare Example insert statement: %v", err)
	}
	insertLookup, err := tx.Prepare("INSERT INTO ExampleLookup VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare ExampleLookup insert statement: %v", err)
	}

	done := 0
	s := bufio.NewScanner(tatoeba)
	// There seem to be duplicate examples in the file, so we dedupe on ID
	seen := make(map[string]struct{})
	var ex Example
	for ex, err = readExample(s); err == nil; ex, err = readExample(s) {
		if _, ok := seen[ex.ID]; ok {
			continue
		}
		seen[ex.ID] = struct{}{}

		if err := convertExample(ex, insertExample, insertLookup); err != nil {
			return fmt.Errorf("could not process example: %v", err)
		}
		done++

		if done%10000 == 0 && progressCB != nil {
			progressCB(done)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("could not read from Tatoeba file: %v", err)
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
	_, err := db.Exec(`CREATE TABLE Example (
		id   TEXT PRIMARY KEY,
		data BLOB NOT NULL     -- Example data in JSON format
	)`)
	if err != nil {
		return fmt.Errorf("could not create Tatoeba example table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE ExampleLookup (
		word       TEXT NOT NULL,
		example_id INTEGER NOT NULL REFERENCES Example(id)
	)`)
	if err != nil {
		return fmt.Errorf("could not create Tatoeba example lookup table: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX ExampleLookup_word ON ExampleLookup(word)`)
	if err != nil {
		return fmt.Errorf("could not create index on Tatoeba example lookup table: %v", err)
	}

	return nil
}

func convertExample(ex Example, insertExample *sql.Stmt, insertLookup *sql.Stmt) error {
	data, err := util.MarshalCompressed(&ex)
	if err != nil {
		return fmt.Errorf("could not marshal example JSON: %v", err)
	}

	if _, err := insertExample.Exec(ex.ID, data); err != nil {
		return fmt.Errorf("could not insert Example data for ID %q: %v", ex.ID, err)
	}
	for _, idx := range ex.Indices {
		if _, err := insertLookup.Exec(idx.Word, ex.ID); err != nil {
			return fmt.Errorf("could not insert ExampleLookup data for ID %q: %v", ex.ID, err)
		}
	}

	return nil
}

// FetchByWord returns all examples using the given word. The results are sorted
// such that "better" examples of the word come first.
func (tb *Tatoeba) FetchByWord(word string) ([]Example, error) {
	rows, err := tb.fetchByWordQuery.Query(word)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var results []Example
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		var result Example
		if err := util.UnmarshalCompressed(data, &result); err != nil {
			return nil, fmt.Errorf("could not unmarshal data: %v", err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].relevance(word) > results[j].relevance(word)
	})
	return results, nil
}

// Example is an example sentence, presented in Japanese and English.
type Example struct {
	ID       string
	Japanese string
	English  string
	Indices  []Index
}

// Segments returns a slice of all the segments in the example.
func (ex Example) Segments() []Segment {
	text := ex.Japanese
	var segs []Segment
	for _, idx := range ex.Indices {
		form := idx.SentenceForm
		if form == "" {
			form = idx.Word
		}

		start := strings.Index(text, form)
		if start < 0 {
			continue
		}
		if start > 0 {
			segs = append(segs, Segment{Text: text[:start]})
		}
		segs = append(segs, Segment{Text: form, Index: &idx})
		text = text[start+len(form):]
	}
	if text != "" {
		segs = append(segs, Segment{Text: text})
	}
	return segs
}

// UniqueIndices returns the indices of the example, removing duplicate entries
// that appear later in the list.
func (ex Example) UniqueIndices() []Index {
	var indices []Index
	seen := make(map[string]struct{})
	for _, idx := range ex.Indices {
		if _, ok := seen[idx.Word]; !ok {
			indices = append(indices, idx)
			seen[idx.Word] = struct{}{}
		}
	}
	return indices
}

var aLineRegexp = regexp.MustCompile(`^A: (.*?)\t(.*?)#ID=([0-9_]+)$`)
var bLineRegexp = regexp.MustCompile(`^B: (.*)$`)

func readExample(s *bufio.Scanner) (Example, error) {
	if !s.Scan() {
		if s.Err() == nil {
			return Example{}, io.EOF
		}
		return Example{}, s.Err()
	}
	aLine := s.Text()
	if !s.Scan() {
		if s.Err() == nil {
			return Example{}, fmt.Errorf("expected B-line after A-line")
		}
		return Example{}, s.Err()
	}
	bLine := s.Text()

	aParts := aLineRegexp.FindStringSubmatch(aLine)
	if aParts == nil {
		return Example{}, fmt.Errorf("unexpected A-line format: %q", aLine)
	}
	bParts := bLineRegexp.FindStringSubmatch(bLine)
	if bParts == nil {
		return Example{}, fmt.Errorf("unexpected B-line format: %q", bLine)
	}

	rawIndices := strings.Fields(bParts[1])
	indices := make([]Index, 0, len(rawIndices))
	for _, raw := range rawIndices {
		idx, err := parseIndex(raw)
		if err != nil {
			return Example{}, err
		}
		indices = append(indices, idx)
	}

	return Example{
		ID:       aParts[3],
		Japanese: aParts[1],
		English:  aParts[2],
		Indices:  indices,
	}, nil
}

// relevance returns a relative "relevance" score of the example to the given
// word. Currently, this is just 0 or 1 depending on whether the example is
// considered a "good" example of the word.
func (ex Example) relevance(word string) int {
	for _, idx := range ex.Indices {
		if word == idx.Word && idx.Good {
			return 1
		}
	}
	return 0
}

// Segment is a segment of Japanese example text which may be associated
// with an Index.
type Segment struct {
	Text  string
	Index *Index
}

// Index is an index for an example sentence, giving details on a word used in
// the sentence.
type Index struct {
	Word           string // the headword as it appears in JMdict
	Disambiguation string // either a reading of the word in kana or an ID of the JMdict entry as #nnnnnnnn
	Sense          string // the number of the sense of the word as used in the sentence
	SentenceForm   string // the form in which the word appears in the sentence
	Good           bool   // whether this sentence is considered a "good example" of the word
}

var indexRegexp = regexp.MustCompile(`^([^[{(~]*)(?:\(([^)]*)\))?(?:\[([^\]]*)\])?(?:\{([^}]*)\})?(~)?`)

func parseIndex(raw string) (Index, error) {
	parts := indexRegexp.FindStringSubmatch(raw)
	if parts == nil {
		return Index{}, fmt.Errorf("unexpected index format: %q", raw)
	}
	return Index{
		Word:           parts[1],
		Disambiguation: parts[2],
		Sense:          parts[3],
		SentenceForm:   parts[4],
		Good:           parts[5] == "~",
	}, nil
}

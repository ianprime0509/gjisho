// Package kanjivg contains types and functions for working with KanjiVG kanji
// drawing data.
package kanjivg

import (
	"bufio"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/ianprime0509/gjisho/internal/util"
)

// KanjiVG is the KanjiVG database, containing data on kanji stroke order and
// drawing.
type KanjiVG struct {
	db         *sql.DB
	fetchQuery *sql.Stmt
}

// New returns a new KanjiVG using the given database.
func New(db *sql.DB) (*KanjiVG, error) {
	fetchQuery, err := db.Prepare("SELECT data FROM StrokeOrder WHERE character = ?")
	if err != nil {
		return nil, fmt.Errorf("could not prepare KanjiVG fetch query: %v", err)
	}
	return &KanjiVG{db, fetchQuery}, nil
}

// ConvertInto converts the KanjiVG data from XML into the given database. The
// given progress callback, if non-nil, is called after every 1,000th converted
// record with the total number of records converted so far.
func ConvertInto(xmlPath string, db *sql.DB, progressCB func(int)) error {
	kanjiVG, err := os.Open(xmlPath)
	if err != nil {
		return fmt.Errorf("could not open KanjiVG file: %v", err)
	}
	defer kanjiVG.Close()

	if err := createTables(db); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	insert, err := tx.Prepare("INSERT INTO StrokeOrder VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare StrokeOrder insert statement: %v", err)
	}

	done := 0
	decoder := xml.NewDecoder(bufio.NewReader(kanjiVG))
	tok, err := decoder.Token()
	for err == nil {
		if start, ok := tok.(xml.StartElement); ok && start.Name.Local == "kanji" {
			if err := convertEntry(decoder, &start, insert); err != nil {
				return fmt.Errorf("could not process KANJIDIC entry: %v", err)
			}
			done++

			if done%1000 == 0 && progressCB != nil {
				progressCB(done)
			}
		}
		tok, err = decoder.Token()
	}
	if err != io.EOF {
		return fmt.Errorf("could not read from KANJIDIC file: %v", err)
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
	_, err := db.Exec(`CREATE TABLE StrokeOrder (
		character TEXT PRIMARY KEY,
		data      BLOB NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("could not create stroke order table: %v", err)
	}

	return nil
}

var kanjiID = regexp.MustCompile("^kvg:kanji_([0-9a-fA-F]{1,5})$")

func convertEntry(decoder *xml.Decoder, start *xml.StartElement, insert *sql.Stmt) error {
	kanji, err := decodeEntry(decoder, start)
	if err != nil {
		return err
	}
	data, err := util.MarshalCompressed(kanji)
	if err != nil {
		return fmt.Errorf("could not marshal stroke order JSON: %v", err)
	}

	_, err = insert.Exec(kanji.Literal, data)
	if err != nil {
		return fmt.Errorf("could not insert stroke order entry: %v", err)
	}
	return nil
}

func decodeEntry(decoder *xml.Decoder, start *xml.StartElement) (Kanji, error) {
	var kanji Kanji
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			if matches := kanjiID.FindStringSubmatch(attr.Value); matches != nil {
				// No error is possible here, because the submatch is guaranteed
				// to be between 1 and 5 hex digits
				cp, _ := strconv.ParseInt(matches[1], 16, 32)
				kanji.Literal = string(rune(cp))
			} else {
				return Kanji{}, fmt.Errorf("invalid kanji ID: %q", attr.Value)
			}
			break
		}
	}
	if kanji.Literal == "" {
		return Kanji{}, errors.New("no ID found for kanji")
	}

	tok, err := decoder.Token()
	for err == nil {
		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "path" {
				for _, attr := range tok.Attr {
					if attr.Name.Local == "d" {
						kanji.Strokes = append(kanji.Strokes, Stroke(attr.Value))
						break
					}
				}
			}
		case xml.EndElement:
			if tok.Name.Local == "kanji" {
				return kanji, nil
			}
		}
		tok, err = decoder.Token()
	}
	return Kanji{}, err
}

// Fetch returns the stroke order data for the given kanji.
func (kvg *KanjiVG) Fetch(kanji string) (Kanji, error) {
	var data []byte
	if err := kvg.fetchQuery.QueryRow(kanji).Scan(&data); err != nil {
		return Kanji{}, fmt.Errorf("scan error: %v", err)
	}

	var character Kanji
	if err := util.UnmarshalCompressed(data, &character); err != nil {
		return Kanji{}, fmt.Errorf("could not unmarshal data: %v", err)
	}

	return character, nil
}

// Kanji is a kanji with associated drawing (stroke) information.
type Kanji struct {
	Literal string
	Strokes []Stroke
}

// Stroke is a stroke of a kanji, represented as an SVG path. See
// https://developer.mozilla.org/en-US/docs/Web/SVG/Tutorial/Paths.
//
// Currently, the only commands supported are M, m, C and c, since those are the
// only ones actually needed for KanjiVG (although I may eventually implement
// the other ones for completeness).
type Stroke string

// DrawTo draws the Stroke to the given Drawer. If markStart is true, a red dot
// is added to the beginning of the stroke to indicate its direction.
func (s Stroke) DrawTo(d Drawer, markStart bool) {
	// We only set the first point after a M or m command
	first := true
	var firstX, firstY float64

	for cmd, args, path := readCommand(string(s)); cmd != ""; cmd, args, path = readCommand(path) {
		switch cmd {
		case "M":
			for len(args) > 0 {
				args = ensureArgs(args, 2)
				d.MoveTo(args[0], args[1])
				if first {
					firstX, firstY = d.GetCurrentPoint()
					first = false
				}
				args = args[2:]
			}
		case "m":
			for len(args) > 0 {
				args = ensureArgs(args, 2)
				x, y := d.GetCurrentPoint()
				d.MoveTo(x+args[0], y+args[1])
				if first {
					firstX, firstY = d.GetCurrentPoint()
					first = false
				}
				args = args[2:]
			}
		case "C":
			for len(args) > 0 {
				args = ensureArgs(args, 6)
				d.CurveTo(args[0], args[1], args[2], args[3], args[4], args[5])
				args = args[6:]
			}
		case "c":
			for len(args) > 0 {
				args = ensureArgs(args, 6)
				x, y := d.GetCurrentPoint()
				d.CurveTo(x+args[0], y+args[1], x+args[2], y+args[3], x+args[4], y+args[5])
				args = args[6:]
			}
		}
	}
	d.Stroke()

	if markStart && !first {
		d.Save()
		d.SetSourceRGB(1, 0, 0)
		r := d.GetLineWidth() * 1.5
		d.MoveTo(firstX, firstY)
		d.Arc(firstX, firstY, r, 0, 2*math.Pi)
		d.Fill()
		d.Restore()
	}
}

var commandRegexp = regexp.MustCompile(`^.*?([MmCc])`)
var argRegexp = regexp.MustCompile(`^[^MmCc]*?(-?[0-9]+\.?[0-9]*)`)

func ensureArgs(args []float64, n int) []float64 {
	for len(args) < n {
		args = append(args, 0)
	}
	return args
}

func readCommand(path string) (cmd string, args []float64, rest string) {
	cmdIdx := commandRegexp.FindStringSubmatchIndex(path)
	if cmdIdx == nil {
		return "", nil, ""
	}
	cmd = path[cmdIdx[2]:cmdIdx[3]]
	path = path[cmdIdx[1]:]

	argIdx := argRegexp.FindStringSubmatchIndex(path)
	for argIdx != nil {
		if arg, err := strconv.ParseFloat(path[argIdx[2]:argIdx[3]], 64); err == nil {
			args = append(args, arg)
		}
		path = path[argIdx[1]:]
		argIdx = argRegexp.FindStringSubmatchIndex(path)
	}

	return cmd, args, path
}

// Drawer is an interface for types that can draw patterns. It is a subset of
// the methods provided by cairo.Context, making it possible to test the drawing
// operations without using Cairo.
type Drawer interface {
	Arc(xc, yc, radius, angle1, angle2 float64)
	CurveTo(x1, y1, x2, y2, x3, y3 float64)
	MoveTo(x, y float64)

	GetCurrentPoint() (x, y float64)
	GetLineWidth() float64
	SetSourceRGB(r, g, b float64)

	Fill()
	Stroke()

	Save()
	Restore()
}

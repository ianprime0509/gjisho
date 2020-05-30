package main

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

// UnmarshalXML unmarshals a NoKanji from XML. This always returns true, since
// the element will be omitted if the value is intended to be false.
func (nk *NoKanji) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if err := d.Skip(); err != nil {
		return err
	}
	*nk = true
	return nil
}

// UnmarshalXMLAttr unmarshals a PartialDescription from an XML attribute.
func (pd *PartialDescription) UnmarshalXMLAttr(a xml.Attr) error {
	*pd = a.Value == "part"
	return nil
}

// UnmarshalXMLAttr unmarshals a Wasei from an XML attribute.
func (w *Wasei) UnmarshalXMLAttr(a xml.Attr) error {
	*w = a.Value == "y"
	return nil
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

func convertDictEntry(decoder *xml.Decoder, start *xml.StartElement, insertEntry *sql.Stmt, insertLookup *sql.Stmt) error {
	var entry DictEntry
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

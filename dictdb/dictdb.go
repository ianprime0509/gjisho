// Package dictdb provides a unified interface to all dictionaries and
// supplemental data stored in GJisho's SQLite database.
package dictdb

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
	"github.com/ianprime0509/gjisho/kanjivg"
	"github.com/ianprime0509/gjisho/kradfile"
	"github.com/ianprime0509/gjisho/tatoeba"
)

// DB is a container for all the more specific dictionaries and databases used
// by GJisho.
type DB struct {
	JMdict   *jmdict.JMdict
	KANJIDIC *kanjidic.KANJIDIC
	KRADFILE *kradfile.KRADFILE
	Tatoeba  *tatoeba.Tatoeba
	KanjiVG  *kanjivg.KanjiVG
}

// Open opens the database, which is looked up using the default path. It
// follows the XDG base directory specification:
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html.
func Open() (*DB, error) {
	dataHome, ok := os.LookupEnv("XDG_DATA_HOME")
	if !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not determine home directory: %v", err)
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	dataDirs, ok := os.LookupEnv("XDG_DATA_DIRS")
	if !ok {
		dataDirs = "/usr/local/share/:/usr/share/"
	}

	lookupDirs := []string{dataHome}
	lookupDirs = append(lookupDirs, strings.Split(dataDirs, ":")...)
	for _, dir := range lookupDirs {
		path := filepath.Join(dir, "gjisho", "gjisho.sqlite")
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return OpenPath(path)
		}
	}

	return nil, fmt.Errorf("database file not found in lookup path %v", lookupDirs)
}

// OpenPath opens the database located at the given path.
func OpenPath(path string) (*DB, error) {
	db := new(DB)

	sqlDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %v", err)
	}

	db.JMdict, err = jmdict.New(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("could not open JMdict handler: %v", err)
	}

	db.KANJIDIC, err = kanjidic.New(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("could not open KANJIDIC handler: %v", err)
	}

	db.KRADFILE, err = kradfile.New(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("could not open KRADFILE handler: %v", err)
	}

	db.Tatoeba, err = tatoeba.New(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("could not open Tatoeba handler: %v", err)
	}

	db.KanjiVG, err = kanjivg.New(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("could not open KanjiVG handler: %v", err)
	}

	return db, nil
}

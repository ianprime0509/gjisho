//go:generate go run github.com/go-bindata/go-bindata/go-bindata -ignore .*~ data/

// Package main contains the main code for GJisho.
package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
	"github.com/ianprime0509/gjisho/kanjivg"
	"github.com/ianprime0509/gjisho/kradfile"
	"github.com/ianprime0509/gjisho/tatoeba"
	_ "github.com/mattn/go-sqlite3"
)

var convertPath = flag.String("conv", "", "convert the given database")
var jmdictPath = flag.String("jmdict", "", "path to the JMdict XML file")
var kanjidicPath = flag.String("kanjidic", "", "path to the KANJIDIC2 XML file")
var kradfilePath = flag.String("kradfile", "", "path to the KRADFILE text file")
var tatoebaPath = flag.String("tatoeba", "", "path to the Tatoeba text file")
var kanjiVGPath = flag.String("kanjivg", "", "path to the KanjiVG XML file")

func main() {
	flag.Parse()
	if *convertPath != "" {
		convert()
	} else {
		LaunchGUI(flag.Args())
	}
}

func convert() {
	db, err := sql.Open("sqlite3", *convertPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	progressCB := func(done int) { log.Printf("Done: %v", done) }
	if *jmdictPath != "" {
		log.Print("Converting JMdict")
		if err := jmdict.ConvertInto(*jmdictPath, db, progressCB); err != nil {
			log.Fatalf("Error converting JMdict: %v", err)
		}
	}
	if *kanjidicPath != "" {
		if err := kanjidic.ConvertInto(*kanjidicPath, db, progressCB); err != nil {
			log.Fatalf("Error converting KANJIDIC: %v", err)
		}
	}
	if *kradfilePath != "" {
		if err := kradfile.ConvertInto(*kradfilePath, db); err != nil {
			log.Fatalf("Error converting KRADFILE: %v", err)
		}
	}
	if *tatoebaPath != "" {
		if err := tatoeba.ConvertInto(*tatoebaPath, db, progressCB); err != nil {
			log.Fatalf("Error converting Tatoeba: %v", err)
		}
	}
	if *kanjiVGPath != "" {
		if err := kanjivg.ConvertInto(*kanjiVGPath, db, progressCB); err != nil {
			log.Fatalf("Error converting KanjiVG: %v", err)
		}
	}
}

package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
	_ "github.com/mattn/go-sqlite3"
)

var convertMode = flag.Bool("conv", false, "convert the given database")
var jmdictPath = flag.String("jmdict", "", "path to the JMdict XML file")
var kanjidicPath = flag.String("kanjidic", "", "path to the Kanjidic2 XML file")

func main() {
	flag.Parse()
	if *convertMode {
		convert(flag.Args())
	} else {
		LaunchGUI(flag.Args())
	}
}

func convert(args []string) {
	db, err := sql.Open("sqlite3", "gjisho.sqlite")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if *jmdictPath != "" {
		if err := jmdict.ConvertInto(*jmdictPath, db); err != nil {
			log.Fatalf("Error converting JMdict: %v", err)
		}
	}
	if *kanjidicPath != "" {
		if err := kanjidic.ConvertInto(*kanjidicPath, db); err != nil {
			log.Fatalf("Error converting Kanjidic: %v", err)
		}
	}
}

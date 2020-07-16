package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
	"github.com/ianprime0509/gjisho/kanjivg"
	"github.com/ianprime0509/gjisho/kradfile"
	"github.com/ianprime0509/gjisho/tatoeba"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: gjisho-cli command [options]")
		os.Exit(2)
	}

	switch cmd := os.Args[1]; cmd {
	case "convert":
		convert(os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, "Unknown command:", cmd)
	}
}

func convert(args []string) {
	fset := flag.NewFlagSet("convert", flag.ExitOnError)
	jmdictPath := fset.String("jmdict", "", "path to the JMdict XML file")
	kanjidicPath := fset.String("kanjidic", "", "path to the KANJIDIC2 XML file")
	kradfilePath := fset.String("kradfile", "", "path to the KRADFILE text file")
	tatoebaPath := fset.String("tatoeba", "", "path to the Tatoeba text file")
	kanjiVGPath := fset.String("kanjivg", "", "path to the KanjiVG XML file")
	fset.Parse(args)
	if fset.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Usage: gjisho-cli convert [flags] db-path")
		os.Exit(2)
	}

	db, err := sql.Open("sqlite3", fset.Arg(0))
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
		log.Print("Converting KANJIDIC")
		if err := kanjidic.ConvertInto(*kanjidicPath, db, progressCB); err != nil {
			log.Fatalf("Error converting KANJIDIC: %v", err)
		}
	}
	if *kradfilePath != "" {
		log.Print("Converting KRADFILE")
		if err := kradfile.ConvertInto(*kradfilePath, db); err != nil {
			log.Fatalf("Error converting KRADFILE: %v", err)
		}
	}
	if *tatoebaPath != "" {
		log.Print("Converting Tatoeba")
		if err := tatoeba.ConvertInto(*tatoebaPath, db, progressCB); err != nil {
			log.Fatalf("Error converting Tatoeba: %v", err)
		}
	}
	if *kanjiVGPath != "" {
		log.Print("Converting KanjiVG")
		if err := kanjivg.ConvertInto(*kanjiVGPath, db, progressCB); err != nil {
			log.Fatalf("Error converting KanjiVG: %v", err)
		}
	}
}

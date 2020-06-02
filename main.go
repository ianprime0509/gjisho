package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/ianprime0509/gjisho/jmdict"
	_ "github.com/mattn/go-sqlite3"
)

var convertMode = flag.Bool("conv", false, "convert the given database")

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

	if err := jmdict.ConvertInto(args[0], db); err != nil {
		log.Fatalf("Error converting JMdict: %v", err)
	}
}

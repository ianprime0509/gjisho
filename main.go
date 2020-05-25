package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	args := os.Args[2:]
	switch os.Args[1] {
	case "convert":
		convert(args)
	case "search":
		search(args)
	}
}

func convert(args []string) {
	if err := ConvertJMdict(args[0], args[1]); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func search(args []string) {
	db, err := OpenJMdict(args[0])
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	entries, err := db.Lookup(args[1])
	if err != nil {
		log.Fatalf("Search error: %v", err)
	}
	fmt.Println(len(entries))
}

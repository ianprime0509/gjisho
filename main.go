package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	args := os.Args[2:]
	switch cmd := os.Args[1]; cmd {
	case "convert":
		convert(args)
	case "launch":
		launch(args)
	case "search":
		search(args)
	default:
		log.Fatalf("Unknown sub-command: %v", cmd)
	}
}

func launch(args []string) {
	LaunchGUI(args)
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
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

package main

import (
	"log"
	"os"

	"github.com/ianprime0509/gjisho/dbconv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := dbconv.ConvertJMdict(os.Args[1], os.Args[2]); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

package main

import (
	"os"

	"github.com/ianprime0509/gjisho/gui"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	gui.LaunchGUI(os.Args)
}

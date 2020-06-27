package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// openDB opens the SQLite database used by GJisho. It follows the XDG base
// directory specification:
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html.
func openDB() (*sql.DB, error) {
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
			return sql.Open("sqlite3", path)
		}
	}

	// Not found; need to try creating the data directory
	dir := filepath.Join(dataHome, "gjisho")
	os.MkdirAll(dir, 0o755)
	return sql.Open("sqlite3", filepath.Join(dir, "gjisho.sqlite"))
}

// removeChildren removes all children from the given container.
func removeChildren(lst *gtk.Container) {
	lst.GetChildren().Foreach(func(item interface{}) {
		lst.Remove(item.(gtk.IWidget))
	})
}

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// OpenDB opens the SQLite database used by GJisho. It follows the XDG base
// directory specification:
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html.
func OpenDB() (*sql.DB, error) {
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

	return nil, fmt.Errorf("database file not found in lookup path %v", lookupDirs)
}

// RemoveChildren removes all children from the given container.
func RemoveChildren(c *gtk.Container) {
	c.GetChildren().Foreach(func(item interface{}) {
		c.Remove(item.(gtk.IWidget))
	})
}

// ScrollToStart scrolls the given scrolled window vertically to the top and
// horizontally to the start.
func ScrollToStart(w *gtk.ScrolledWindow) {
	w.GetVAdjustment().SetValue(w.GetVAdjustment().GetLower())
	w.GetHAdjustment().SetValue(w.GetHAdjustment().GetLower())
}

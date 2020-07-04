package util

import (
	"bufio"
	"bytes"
	"compress/flate"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// MarshalCompressed marshals the given value into compressed JSON.
func MarshalCompressed(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("JSON error: %v", err)
	}
	buf := new(bytes.Buffer)
	w, err := flate.NewWriter(buf, flate.DefaultCompression)
	if err != nil {
		return nil, fmt.Errorf("could not create compressor: %v", err)
	}
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("could not compress data: %v", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("could not compress data: %v", err)
	}
	return buf.Bytes(), nil
}

// UnmarshalCompressed unmarshals the given value from compressed JSON.
func UnmarshalCompressed(b []byte, v interface{}) error {
	r := flate.NewReader(bytes.NewReader(b))
	data := new(bytes.Buffer)
	if _, err := io.Copy(data, r); err != nil {
		return fmt.Errorf("could not decompress data: %v", err)
	}
	if err := r.Close(); err != nil {
		return fmt.Errorf("could not decompress data: %v", err)
	}
	if err := json.Unmarshal(data.Bytes(), v); err != nil {
		return fmt.Errorf("JSON error: %v", err)
	}
	return nil
}

var entityRegexp = regexp.MustCompile(`<!ENTITY +([^" ]+) *"([^"]*)">`)
var endDoctypeRegexp = regexp.MustCompile(`]>`)

// ParseEntities parses the entity definitions from the doctype of the given XML
// file. It is not very smart and does only enough to be useful.
func ParseEntities(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()
	input := bufio.NewReader(file)

	entities := make(map[string]string)
	line, err := input.ReadString('\n')
	for err == nil {
		if endDoctypeRegexp.MatchString(line) {
			return entities, nil
		}
		if matches := entityRegexp.FindStringSubmatch(line); matches != nil {
			entities[matches[1]] = matches[2]
		}

		line, err = input.ReadString('\n')
	}

	if err != io.EOF {
		return nil, fmt.Errorf("could not read from file: %v", err)
	}
	return entities, nil
}

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

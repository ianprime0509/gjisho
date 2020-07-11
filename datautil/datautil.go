package datautil

import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
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

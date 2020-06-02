package xmlutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

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

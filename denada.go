package denada

import "os"
import "io"
import "fmt"
import "strings"

// This file contains the API for the denada parser

var errorList []error

func listToError(l []error) error {
	msg := "Parsing errors:"
	for _, e := range l {
		msg += fmt.Sprintf("\n  %v", e)
	}
	return fmt.Errorf("%s", msg)
}

func ParseString(s string) (ElementList, error) {
	r := strings.NewReader(s)
	return Parse(r)
}

func ParseFile(filename string) (ElementList, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	p, err := NewParser(r, filename)
	if err != nil {
		return nil, err
	}
	return p.ParseFile()
}

func Parse(r io.Reader) (ElementList, error) {
	p, err := NewParser(r, "")
	if err != nil {
		return nil, err
	}
	return p.ParseFile()
}

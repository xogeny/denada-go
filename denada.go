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

func (yylex Lexer) Error(e string) {
	errorList = append(errorList,
		fmt.Errorf("Error %s at line %d, column %d", e, lineNumber, colNumber))
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

	return Parse(r)
}

func Parse(r io.Reader) (ElementList, error) {
	errorList = []error{}
	lineNumber = 0
	colNumber = 0

	lex := NewLexer(r)
	ret := yyParse(lex)
	if ret == 0 {
		return _parserResult, nil
	} else {
		ret := []error{}
		for _, msg := range errorList {
			ret = append(ret, msg)
		}
		return nil, listToError(ret)
	}
}

package denada

import "io"
import "fmt"
import "strings"

// This file contains the API for the denada parser

var errorList []error

func (yylex Lexer) Error(e string) {
	errorList = append(errorList,
		fmt.Errorf("Error %s at line %d, column %d", e, lineNumber, colNumber))
}

func ParseString(s string) (ElementList, []error, bool) {
	r := strings.NewReader(s)
	return Parse(r)
}

func Parse(r io.Reader) (ElementList, []error, bool) {
	errorList = []error{}
	lineNumber = 0
	colNumber = 0

	lex := NewLexer(r)
	ret := yyParse(lex)
	if ret == 0 {
		return _parserResult, []error{}, true
	} else {
		ret := []error{}
		for _, msg := range errorList {
			ret = append(ret, msg)
		}
		return nil, ret, false
	}
}

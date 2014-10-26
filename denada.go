package denada

import "io"
import "fmt"

// This file contains the API for the denada parser

func Parse(r io.Reader) (ElementList, error) {
	lex := NewLexer(r)
	ret := yyParse(lex)
	if ret == 0 {
		return _parserResult, nil
	} else {
		return nil, fmt.Errorf("Parsing error")
	}
}

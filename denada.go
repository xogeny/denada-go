package denada

import "io"
import "fmt"

// This file contains the API for the denada parser

var errorList []error

func (yylex Lexer) Error(e string) {
	errorList = append(errorList,
		fmt.Errorf("Error %s at line %d, column %d", e, yylex.l, yylex.c))
}

func Parse(r io.Reader) (ElementList, error) {
	errorList = []error{}
	lex := NewLexer(r)
	ret := yyParse(lex)
	if ret == 0 {
		return _parserResult, nil
	} else {
		msg := ""
		for _, emsg := range errorList {
			msg = msg + fmt.Sprintf("%v\n", emsg)
		}
		return nil, fmt.Errorf(msg)
	}
}

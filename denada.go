package denada

import "io"

// This file contains the API for the denada parser

func Parse(r io.Reader) error {
	lex := NewLexer(r)
	_ = yyParse(lex)
	return nil
}

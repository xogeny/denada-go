package main

import "io"
import "fmt"
import "text/scanner"

type ParseError struct {
	Expected string
	Found string
	Position scanner.Position
}

func MakeError(s scanner.Scanner, msg string) ParseError {
	return ParseError{Position: s.Pos(), Expected: msg, Found: s.TokenText()};
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("Expected %s but found '%s' at line %d, column %d", pe.Expected,
		pe.Found, pe.Position.Line, pe.Position.Column);
}

func Tokenize(src io.Reader) {
	var s scanner.Scanner
	s.Init(src)
	tok := s.Scan()
	for tok != scanner.EOF {
		fmt.Printf("Token: '%s' %v\n", s.TokenText(), tok);
		// do something with tok
		tok = s.Scan()
	}
}

func Parse(src io.Reader) error {
	var s scanner.Scanner
	s.Init(src)

	return ParseFile(s);
}

func ParseFile(s scanner.Scanner) error {
	ids := []string{};
	tok := s.Peek();

	if (tok==scanner.EOF) { return nil; }
	if (tok=='}') { return nil; }

	tok = s.Scan();

	for ; tok==scanner.Ident;  {
		ids = append(ids, s.TokenText());
		tok = s.Scan();
	}

	switch tok {
	case '=':
		fmt.Println("Declaration");
		return nil;
	case ';':
		fmt.Println("Declaration");
		return nil;
	case '{':
		fmt.Println("Definition");
		ParseFile(s);
		tok = s.Scan();
		if (tok!='{') { return MakeError(s, "'{'"); }
		return nil;
	default:
		return MakeError(s, "'=', ';' or '{'");
	}
}

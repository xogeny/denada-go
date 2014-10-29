package denada

import "fmt"

type TokenType int

const (
	T_IDENTIFIER TokenType = iota
	T_RBRACE
	T_LBRACE
	T_LPAREN
	T_RPAREN
	T_QUOTE
	T_EQUALS
	T_SEMI
	T_COMMA
	T_EOF
	T_WHITE
	T_UNKNOWN
)

func (tt TokenType) String() string {
	switch tt {
	case T_IDENTIFIER:
		return "<identifier>"
	case T_RBRACE:
		return "}"
	case T_LBRACE:
		return "{"
	case T_LPAREN:
		return "("
	case T_RPAREN:
		return ")"
	case T_QUOTE:
		return "\""
	case T_EQUALS:
		return "="
	case T_SEMI:
		return ";"
	case T_COMMA:
		return ","
	case T_WHITE:
		return "<whitespace>"
	case T_EOF:
		return "EOF"
	case T_UNKNOWN:
		fallthrough
	default:
		return "<???>"
	}
}

type UnexpectedToken struct {
	Found    Token
	Expected string
}

func (u UnexpectedToken) Error() string {
	if u.Found.File != "" {
		return fmt.Sprintf("Expecting %s, found '%v' @ (%d, %d) in %s", u.Expected, u.Found.Type,
			u.Found.Line, u.Found.Column, u.Found.File)
	} else {
		return fmt.Sprintf("Expecting %s, found '%v' @ (%d, %d)", u.Expected, u.Found.Type,
			u.Found.Line, u.Found.Column)
	}
}

type Token struct {
	Type   TokenType
	String string
	Line   int
	Column int
	File   string
}

func (t Token) Expected(expected string) UnexpectedToken {
	return UnexpectedToken{
		Found:    t,
		Expected: expected,
	}
}

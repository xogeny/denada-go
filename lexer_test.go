package denada

import "testing"
import "strings"

import . "github.com/onsi/gomega"

func Test_LexerStream(t *testing.T) {
	RegisterTestingT(t)

	var lval yySymType
	r := strings.NewReader("set x = 5;")
	lexer := NewLexer(r)
	tok1 := lexer.Lex(&lval)

	Expect(tok1).To(Equal(IDENTIFIER))
	Expect(lval.identifier).To(Equal("set"))

	tok2 := lexer.Lex(&lval)

	Expect(tok2).To(Equal(IDENTIFIER))
	Expect(lval.identifier).To(Equal("x"))

	tok3 := lexer.Lex(&lval)

	Expect(tok3).To(Equal(EQUALS))

	tok4 := lexer.Lex(&lval)

	Expect(tok4).To(Equal(NUMBER))
	Expect(lval.number).To(Equal(5))

	tok5 := lexer.Lex(&lval)

	Expect(tok5).To(Equal(SEMI))
}

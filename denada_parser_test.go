package denada

import "testing"
import "strings"

import . "github.com/onsi/gomega"

var sample = `
printer 'ABC' {
   set location = "Mike's desk";
   set model = "HP 8860";
}

'printer' DEF {
   set location = "Coffee machine";
   set model = "HP 8860";
   set networkName = "PrinterDEF";
}

computer XYZ {
   set location = "Mike's desk";
   set 'model' = "Mac Book Air";
}
`

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

func Test_SimpleDeclaration(t *testing.T) {
	RegisterTestingT(t)

	r := strings.NewReader("set x = 5 \"Description\";")
	elems, err := Parse(r)

	Expect(err).To(BeNil())
	Expect(len(elems)).To(Equal(1))

	elem := elems[0]

	Expect(elem.isDeclaration()).To(BeTrue())
	Expect(elem.isDefinition()).To(BeFalse())
	Expect(len(elem.Modifications)).To(Equal(0))

	Expect(elem.Qualifiers).To(Equal([]string{"set"}))
	Expect(elem.Name).To(Equal("x"))
	Expect(elem.Description).To(Equal("Description"))
	Expect(elem.Value).To(Equal(5))
}

func Test_SampleInput(t *testing.T) {
	RegisterTestingT(t)

	r := strings.NewReader(sample)
	elems, err := Parse(r)
	Expect(err).To(BeNil())
	Expect(len(elems)).To(Equal(3))
}

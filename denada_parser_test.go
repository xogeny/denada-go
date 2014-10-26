package denada

import "testing"
import "log"
import "strings"

import . "github.com/xogeny/gocore/xassert"

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
	ExtraLayer(t, "Testing lexer", func() {
		var lval yySymType
		r := strings.NewReader("set x = 5;")
		lexer := NewLexer(r)
		tok1 := lexer.Lex(&lval)
		IsEqual(tok1, IDENTIFIER)
		IsEqual(lval.identifier, "set")
		tok2 := lexer.Lex(&lval)
		IsEqual(tok2, IDENTIFIER)
		IsEqual(lval.identifier, "x")
		tok3 := lexer.Lex(&lval)
		IsEqual(tok3, EQUALS)
		tok4 := lexer.Lex(&lval)
		IsEqual(tok4, NUMBER)
		tok5 := lexer.Lex(&lval)
		IsEqual(tok5, SEMI)
	})
}

func Test_SimpleDeclaration(t *testing.T) {
	ExtraLayer(t, "Testing parser", func() {
		r := strings.NewReader("set x = 5 \"Description\";")
		elems, err := Parse(r)
		NoError(err)
		IsEqual(1, len(elems))
		elem := elems[0]
		log.Printf("Elem = %v", elem)
		IsTrue(elem.isDeclaration())
		IsFalse(elem.isDefinition())
		IsEqual(0, len(elem.Modifications))
		Resembles(elem.Qualifiers, []string{"set"})
		IsEqual(elem.Name, "x")
		IsEqual(elem.Description, "Description")
		Resembles(elem.Value, 5)
	})
}

func Test_SampleInput(t *testing.T) {
	ExtraLayer(t, "Testing parser", func() {
		r := strings.NewReader(sample)
		elems, err := Parse(r)
		NoError(err)
		IsEqual(3, len(elems))
	})
}

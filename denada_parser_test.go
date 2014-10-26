package denada

import "testing"
import "log"
import "strings"

import . "github.com/xogeny/gocore/xassert"

var sample = `
printer ABC {
   set location = "Mike's desk";
   set model = "HP 8860";
}

printer DEF {
   set location = "Coffee machine";
   set model = "HP 8860";
   set networkName = "PrinterDEF";
}

computer XYZ {
   set location = "Mike's desk";
   set model = "Mac Book Air";
}
`

func Test1(t *testing.T) {
	ExtraLayer(t, "Testing parser", func() {
		r := strings.NewReader("set x = 5;")
		elems, err := Parse(r)
		if err != nil {
			log.Printf("Error: %x", err)
		}
		IsEqual(1, len(elems))
		elem := elems[0]
		log.Printf("Elem = %v", elem)
		Resembles(elem.Qualifiers, []string{"set"})
		IsEqual(elem.Name, "x")
		Resembles(elem.Value, 5)
	})
}

func Test2(t *testing.T) {
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

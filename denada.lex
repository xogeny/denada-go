/;/               { yylex.move(); return SEMI; }
/\(/              { yylex.move(); return LPAREN; }
/\)/              { yylex.move(); return RPAREN; }
/\{/              { yylex.move(); return LBRACE; }
/\}/              { yylex.move(); return RBRACE; }
/=/               { yylex.move(); return EQUALS; }
/,/               { yylex.move(); return COMMA; }
/true/            { yylex.move(); lval.bool = true; return BOOLEAN; }
/false/           { yylex.move(); lval.bool = false; return BOOLEAN; }
/[0-9][0-9]*/     { yylex.move(); lval.number,_ = strconv.Atoi(yylex.Text()); return NUMBER }
/'[^']*'/         {
   yylex.move(); lval.identifier = strings.Trim(yylex.Text(), "'"); return IDENTIFIER
}
/"[^"]*"/         { yylex.move(); lval.string = strings.Trim(yylex.Text(), "\""); return STRING }
/[A-Za-z_][A-Za-z_0-9]*/ { yylex.move(); lval.identifier = yylex.Text(); return IDENTIFIER; }
/[ \t]+/          { yylex.move(); /* eat up whitespace */ }
/[\n]+/           { yylex.move(); }
/./               { yylex.move(); println("Unrecognized character:", yylex.Text()) }
/{[^\{\}\n]*}/    { yylex.move(); /* eat up one-line comments */ }
//
package denada

import "fmt"
import "strconv"

var lineNumber int;
var colNumber int;

func (yylex Lexer) move() {
    for _, c := range(yylex.Text()) {
		if c=='\n' {
			colNumber = 0;
			lineNumber++;
  	    } else {
		    // TODO: Need to handle non-printable characters (e.g. \r)
		    colNumber++;
		}
	}
}

func ystream(r io.Reader) {
  lexer := NewLexer(r);
  for {
	var lval yySymType;
    tok := lexer.Lex(&lval);
	if tok==0 {
		break
	}
	fmt.Printf("Token #%d '%s' (lval=%v)\n", tok, lexer.Text(), lval);
  }
}

func stream(r io.Reader) {
  lexer := NewLexer(r);
  for {
    tok := lexer.next(0);
	if tok==-1 {
		break
	}
	fmt.Printf("Token #%d '%s'\n", tok, lexer.Text());
  }
}

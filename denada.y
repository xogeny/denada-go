%{

  /*

	Copyright 2014: Xogeny, Inc.

	CAUTION: If this file is a Go source file (*.go), it was generated
	automatically by '$ go tool yacc' from a *.y file - DO NOT EDIT in that case!

  */

package denada

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cznic/strutil"
)

%}

%union {
    identifier string
	string     string
	bool       bool
	item interface{}
}

%token	BOOLEAN
%token	IDENTIFIER
%token	NUMBER
%token	STRING

%type <bool> BOOLEAN
%type <item> IDENTIFIER
%type <item> NUMBER
%type <string> STRING

%type	<item> 	/*TODO real type(s), if/where applicable */
	Declaration
	Declaration2
	Definition
	Definition1
	Expr
	File
	File1
	File11
	Modification
	Modifiers
	Modifiers1
	Modifiers11
	Preface
	Preface1
	QualifiersAndId
	QualifiersAndId1
	Start

%start Start

%%

Declaration
: Preface Declaration2 ';' { $$ = []Declaration{$1, $2, ";"} /* TODO 1 */ }
| Preface '=' Expr Declaration2 ';' { $$ = []Declaration{$1, "=", $3, $4, ";"} /* TODO 2 */	}

Declaration2
: /* EMPTY */ { $$ = nil /* TODO 3 */ }
| STRING { $$ = $1  /* TODO 4 */	}

Definition
: Preface Definition1 '{' File '}' { $$ = []Definition{$1, $2, "{", $4, "}"} /* TODO 5 */ }

Definition1
: /* EMPTY */ {	$$ = nil /* TODO 6 */ }
| STRING { $$ = $1 /* TODO 7 */	}

Expr
: STRING { $$ = $1 /* TODO 8 */ }
| NUMBER { $$ = $1 /* TODO 9 */ }
| BOOLEAN { $$ = $1 /* TODO 10 */ }

File
: File1 { $$ = $1 /* TODO 11 */ }

File1
: /* EMPTY */ {	$$ = []File1(nil) /* TODO 12 */ }
| File1 File11 { $$ = append($1.([]File1), $2) /* TODO 13 */ }

File11
: Definition { $$ = $1 /* TODO 14 */ }
| Declaration {	$$ = $1 /* TODO 15 */ }

Modification
: IDENTIFIER '=' Expr { $$ = []Modification{$1, "=", $3} /* TODO 16 */ }

Modifiers
: '(' Modifiers1 ')' { $$ = []Modifiers{"(", $2, ")"} /* TODO 17 */ }

Modifiers1
: /* EMPTY */ {	$$ = nil /* TODO 18 */ }
| Modification Modifiers11 { $$ = []Modifiers1{$1, $2} /* TODO 19 */ }

Modifiers11
: /* EMPTY */ {	$$ = []Modifiers11(nil) /* TODO 20 */ }
| Modifiers11 ',' Modification { $$ = append($1.([]Modifiers11), ",", $3) /* TODO 21 */ }

Preface
: QualifiersAndId Preface1 { $$ = []Preface{$1, $2} /* TODO 22 */ }

Preface1
: /* EMPTY */ {	$$ = nil /* TODO 23 */ }
| Modifiers	{ $$ = $1 /* TODO 24 */ }

QualifiersAndId
: QualifiersAndId1 IDENTIFIER {	$$ = []QualifiersAndId{$1, $2} /* TODO 25 */ }

QualifiersAndId1
: /* EMPTY */ {	$$ = []QualifiersAndId1(nil) /* TODO 26 */ }
| QualifiersAndId1 IDENTIFIER { $$ = append($1.([]QualifiersAndId1), $2) /* TODO 27 */ }

Start
: File { _parserResult = $1 /* TODO 28 */ }

%%

//TODO remove demo stuff below

var _parserResult interface{}

type (
	Declaration interface{}
	Declaration2 interface{}
	Definition interface{}
	Definition1 interface{}
	File interface{}
	File1 interface{}
	File11 interface{}
	Modification interface{}
	Modifiers interface{}
	Modifiers1 interface{}
	Modifiers11 interface{}
	Preface interface{}
	Preface1 interface{}
	QualifiersAndId interface{}
	QualifiersAndId1 interface{}
	Start interface{}
)

func _dump() {
	s := fmt.Sprintf("%#v", _parserResult)
	s = strings.Replace(s, "%", "%%", -1)
	s = strings.Replace(s, "{", "{%i\n", -1)
	s = strings.Replace(s, "}", "%u\n}", -1)
	s = strings.Replace(s, ", ", ",\n", -1)
	var buf bytes.Buffer
	strutil.IndentFormatter(&buf, ". ").Format(s)
	buf.WriteString("\n")
	a := strings.Split(buf.String(), "\n")
	for _, v := range a {
		if strings.HasSuffix(v, "(nil)") || strings.HasSuffix(v, "(nil),") {
			continue
		}
	
		fmt.Println(v)
	}
}

// End of demo stuff

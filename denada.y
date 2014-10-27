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
    bool       bool
    number     interface{}
	string     string
	elements   ElementList
	element    *Element
    expr       Expr
    dict       map[string]Expr
}

%token	BOOLEAN
%token	IDENTIFIER
%token	NUMBER
%token	STRING
%token  SEMI
%token  LPAREN
%token  RPAREN
%token  LBRACE
%token  RBRACE
%token  EQUALS
%token  COMMA

// Tokens
%type <bool> BOOLEAN
%type <identifier> IDENTIFIER
%type <number> NUMBER
%type <string> STRING

// Rules
%type <expr> Expr

%type <string> Description

%type <element> Elem
%type <element> Declaration
%type <element> Definition
%type <element> Preface
%type <element> QualifiersAndId QualifiersAndId1

%type <elements> File File1
%type <elements> Start

%type <dict> Modification Modifiers Modifiers1 Modifiers11 PrefaceModifiers

%start Start

%%

Declaration
: Preface Description SEMI {
  $$ = $1;
  $$.Description = $2;
}
| Preface EQUALS Expr Description SEMI {
  $$ = $1;
  $$.Value = $3;
  $$.Description = $4;
}

Description
: /* EMPTY */ { $$ = "" }
| STRING { $$ = $1 }

Definition
: Preface Description LBRACE File RBRACE {
  $$ = $1;
  $$.Description = $2;
  $$.Contents = $4;
}

Expr
: STRING { $$ = $1 }
| NUMBER { $$ = $1 }
| BOOLEAN { $$ = $1 }
| IDENTIFIER { $$ = $1 }

File
: File1 { $$ = $1 }

File1
: /* EMPTY */ {	$$ = MakeElementList() }
| File1 Elem {
  $$ = append($1, $2);
}

Elem
: Definition { $$ = $1; $$.definition = true; }
| Declaration {	$$ = $1; $$.definition = false; }

Modification
: IDENTIFIER EQUALS Expr {
  $$ = map[string]Expr{};
  $$[$1] = $3;
}

Modifiers
: LPAREN Modifiers1 RPAREN {
  $$ = $2;
}

Modifiers1
: /* EMPTY */ {	$$ = map[string]Expr{} }
| Modification Modifiers11 {
  $$ = $2;
  for k, v := range($1) {
	  $$[k] = v;
  }
}

Modifiers11
: /* EMPTY */ {	$$ = map[string]Expr{} }
| Modifiers11 COMMA Modification {
  $$ = $1;
  for k, v := range($3) {
	  $$[k] = v;
  }
}

Preface
: QualifiersAndId PrefaceModifiers {
  $$ = $1;
  $1.Modifications = $2;
}

PrefaceModifiers
: /* EMPTY */ {	$$ = map[string]Expr{} }
| Modifiers	{ $$ = $1 }

QualifiersAndId
: QualifiersAndId1 IDENTIFIER {
  if $1.Name!="" {
    $1.Qualifiers = append($1.Qualifiers, $1.Name);
  }
  $1.Name = $2;
  $$ = $1;
}

QualifiersAndId1
: /* EMPTY */ {	$$ = &Element{} }
| QualifiersAndId1 IDENTIFIER {
  if $1.Name!="" {
    $1.Qualifiers = append($1.Qualifiers, $1.Name);
  }
  $1.Name = $2;
  $$ = $1;
}

Start
: File {
  _parserResult = $1;
}

%%

var _parserResult ElementList

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

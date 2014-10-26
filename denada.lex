/;/               { return SEMI; }
/\(/              { return LPAREN; }
/\)/              { return RPAREN; }
/\{/              { return LBRACE; }
/\}/              { return RBRACE; }
/=/               { return EQUALS; }
/,/               { return COMMA; }
/true/            { lval.bool = true; return BOOLEAN; }
/false/           { lval.bool = false; return BOOLEAN; }
/[A-Za-z_][A-Za-z_0-9]*/ { lval.identifier = yylex.Text(); return IDENTIFIER; }
/[ \t\n]+/        { /* eat up whitespace */ }
/./               { println("Unrecognized character:", yylex.Text()) }
/{[^\{\}\n]*}/    { /* eat up one-line comments */ }
//
package denada

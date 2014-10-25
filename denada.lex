/[A-Za-z_][A-Za-z_0-9]+/ { lval.identifier = yylex.Text(); return IDENTIFIER }
/[ \t\n]+/        { /* eat up whitespace */ }
/./               { println("Unrecognized character:", yylex.Text()) }
/{[^\{\}\n]*}/    { /* eat up one-line comments */ }
//
package denada

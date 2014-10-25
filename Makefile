scanner.go: denada.l
	golex -o $@ $<

parser.go: denada.y
	go tool yacc -o $@ $<

denada.y: denada.ebnf
	ebnf2y -start File -o denada.y denada.ebnf

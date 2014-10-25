scanner.go: denada.l
	golex -o $@ $<

parser.go: denada.y
	go tool yacc -v y.output -o $@ $<

denada.y: denada.ebnf
	ebnf2y -start File -m -o denada.y denada.ebnf

all: denada_scanner.go denada_parser.go

denada_scanner.go: denada.l
	golex -o $@ $<
	go fmt $@

denada_parser.go: denada.y
	go tool yacc -v y.output -o $@ $<
	go fmt $@

denada.y: denada.ebnf
	ebnf2y -start File -m -o denada.y denada.ebnf

clean:
	-rm y.go y.output denada.y denada_parser.go denada_scanner.go

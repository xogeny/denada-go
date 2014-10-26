all: denada_scanner.go denada_parser.go

#denada_scanner.go: denada.l
#	golex -o $@ $<
#	go fmt $@

denada_scanner.go: denada.lex
	nex -e -o $@ $<
	go fmt $@

denada_parser.go: denada.y
	go tool yacc -v y.output -o $@ $<
	go fmt $@

#denada.y: denada.ebnf
#	ebnf2y -pkg denada -start File -M -o denada.y denada.ebnf
#	-rm y.go

clean:
	-rm y.go y.output denada_parser.go denada_scanner.go

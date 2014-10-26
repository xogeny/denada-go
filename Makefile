all: denada_scanner.go denada_parser.go

denada_scanner.go: denada.lex
	nex -e -o $@ $<
	go fmt $@

denada_parser.go: denada.y
	go tool yacc -v y.output -o $@ $<
	go fmt $@

clean:
	-rm y.go y.output denada_parser.go denada_scanner.go

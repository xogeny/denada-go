package main

import "os"

import "github.com/jessevdk/go-flags"

func main() {
	var options struct{}

	parser := flags.NewParser(&options, flags.Default)

	parser.AddCommand("format",
		"Rewrite a Denada file in canonical form",
		"Rewrite a Denada file in canonical form",
		&FormatCommand{})

	parser.AddCommand("parse",
		"Parse a Denada file",
		"Parse a Denada file",
		&ParseCommand{})

	parser.AddCommand("check",
		"Parse a Denada file and check it against a grammar file",
		"Parse a Denada file and check it against a grammar file",
		&CheckCommand{})

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}

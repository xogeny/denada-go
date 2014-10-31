package main

import "fmt"
import "github.com/xogeny/denada-go"

type CheckCommand struct {
	Positional struct {
		Input   string `description:"Input file"`
		Grammar string `description:"Grammar file"`
	} `positional-args:"true" required:"true"`
}

func (f CheckCommand) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments")
	}

	ifile := f.Positional.Input
	elems, err := denada.ParseFile(ifile)
	if err != nil {
		return fmt.Errorf("Error parsing input file %s: %v", ifile, err)
	}

	gfile := f.Positional.Grammar
	grammar, err := denada.ParseFile(gfile)
	if err != nil {
		return fmt.Errorf("Error parsing grammar file %s: %v", gfile, err)
	}

	err = denada.Check(elems, grammar, false)
	if err != nil {
		denada.Check(elems, grammar, true)
		return fmt.Errorf("File %s was not a valid instance of the grammar in %s",
			ifile, gfile)
	}

	return nil
}

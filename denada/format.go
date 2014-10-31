package main

import "os"
import "fmt"
import "github.com/xogeny/denada-go"

type FormatCommand struct {
	Positional struct {
		Term string `description:"Input file"`
	} `positional-args:"true" required:"true"`
}

func (f FormatCommand) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments")
	}

	file := f.Positional.Term
	elems, err := denada.ParseFile(file)
	if err != nil {
		return fmt.Errorf("Error parsing input file %s: %v", file, err)
	}

	fp, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Error rewriting %s: %v", file, err)
	}
	defer fp.Close()
	denada.UnparseTo(elems, fp)
	return nil
}

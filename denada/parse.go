package main

import "os"
import "fmt"
import "path"
import "github.com/xogeny/denada-go"

type ParseCommand struct {
	Positional struct {
		Input string `description:"Input file"`
	} `positional-args:"true" required:"true"`
	Import bool `short:"i" long:"import" description:"Expand imports"`
	Echo   bool `short:"e" long:"echo" description:"Echo parsed data"`
}

func (f ParseCommand) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments")
	}

	file := f.Positional.Input

	err := os.Chdir(path.Dir(file))
	if err != nil {
		return fmt.Errorf("Error changing directory to %s: %v", path.Dir(file), err)
	}

	elems, err := denada.ParseFile(file)
	if err != nil {
		return fmt.Errorf("Error parsing input file %s: %v", file, err)
	}

	if f.Import {
		elems, err = denada.ImportTransform(elems)
		if err != nil {
			return fmt.Errorf("Error doing imports in %s: %v", file, err)
		}
	}

	if f.Echo {
		denada.UnparseTo(elems, os.Stdout)
	}

	fmt.Printf("File %s is syntactically correct Denada\n", file)
	return nil
}

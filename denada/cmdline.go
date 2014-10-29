package main

import "github.com/xogeny/denada-go"
import "os"
import "fmt"

func main() {
	if len(os.Args) == 2 {
		_, err := denada.ParseFile(os.Args[1])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("Successfully parsed\n")
		os.Exit(0)
	} else if len(os.Args) == 3 {
		input, err := denada.ParseFile(os.Args[1])
		if err != nil {
			fmt.Printf("Error parsing input file: %v\n", err)
			os.Exit(2)
		}
		grammar, err := denada.ParseFile(os.Args[2])
		if err != nil {
			fmt.Printf("Error parsing grammar file: %v\n", err)
			os.Exit(3)
		}
		err = denada.Check(input, grammar, false)
		if err != nil {
			// Run the check again, but this time turn on diagnostics
			denada.Check(input, grammar, true)
			fmt.Printf("Error checking input against grammar: %v\n", err)
			os.Exit(4)
		}
		fmt.Printf("Successfully parsed and checked\n")
		os.Exit(0)
	} else {
		fmt.Printf("Usage: %s input [grammar]\n", os.Args[0])
		os.Exit(1)
	}
}

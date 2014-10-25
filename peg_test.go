package main

import "log"
import "testing"

import . "github.com/onsi/gomega"

func TestPEG(t *testing.T) {
	RegisterTestingT(t);

	log.Printf("Creating parser");
	d := Denada{};
	d.Buffer = "class X { x; }";
	d.Ids = []string{};

	log.Printf("Initializing parser");
	d.Init();

	log.Printf("Parsing...");
	err := d.Parse();
	if (err!=nil) {
		log.Printf("Error: %s", err.Error());
	}
	Expect(err).To(BeNil());
	log.Printf("...done");

	d.Execute();

	//d.tokenTree.AST().Print(d.Buffer);

	log.Printf("d = %v", d);
}

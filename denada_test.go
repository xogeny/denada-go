package main

import "os"
import "log"
import "io/ioutil"
import "testing"

import . "github.com/onsi/gomega"

import "github.com/robertkrimen/otto"

func testSimple(t *testing.T) {
	RegisterTestingT(t);

	sample := `
printer ABC {
   set location = "Mike's desk";
   set model = "HP 8860";
}

printer DEF {
   set location = "Coffee machine";
   set model = "HP 8860";
   set networkName = "PrinterDEF";
}

computer XYZ {
   set location = "Mike's desk";
   set model = "Mac Book Air";
}
`

	vm := otto.New();

	gf, err := os.Open("grammar.js");
	Expect(err).To(BeNil());

	raw, err := ioutil.ReadAll(gf);
	Expect(err).To(BeNil());

	_, err = vm.Run(raw);
	Expect(err).To(BeNil());

	denada, err := vm.Object("denada");
	Expect(err).To(BeNil());
	
	input, err := vm.ToValue(sample);
	Expect(err).To(BeNil());

	ast, err := denada.Call("parse", input);
	Expect(err).To(BeNil());

	log.Printf("AST = %v", ast);
}

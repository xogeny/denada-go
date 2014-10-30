package denada

import "testing"
import . "github.com/onsi/gomega"

var importTest = `
import(file="testsuite/case1.dnd", recursive=false);

import(file="testsuite/case1.dnd");

scoped {
  import(file="testsuite/case2.dnd", recursive=true);
}
`

var expectedResult = `
// Most basic syntactic example with just 2 declarations

props(declarations=2);

Real r;

// Most basic syntactic example with just 2 declarations

props(declarations=2);

Real r;

scoped {
  // Checking against a grammar file with different expression types

  props(grammar="config.grm", definitions=1, declarations=1) "props";

  section Foo "section" {
     x = 1 "variable";
     y = 1.0 "variable";
     z = "test string" "variable";
     json = {"this": "is a JSON expression!"} "variable";
  }
}
`

func Test_NonRecursiveImport(t *testing.T) {
	RegisterTestingT(t)

	exp, err := ParseString(expectedResult)
	Expect(err).To(BeNil())

	raw, err := ParseString(importTest)
	Expect(err).To(BeNil())

	elab, err := ImportTransform(raw)
	Expect(err).To(BeNil())

	eq := elab.Equals(exp)
	Expect(eq).To(BeNil())
}

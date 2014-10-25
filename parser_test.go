package main

import "testing"
import "strings"

import . "github.com/onsi/gomega"

func testParser(t *testing.T) {
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

	buf := strings.NewReader(sample);
	err := Parse(buf);
	Expect(err).To(BeNil());
}

package denada

import "testing"
import "strings"

import . "github.com/onsi/gomega"

var sample = `
printer 'ABC' {
   set location = "Mike's desk";
   set model = "HP 8860";
}

'printer' DEF {
   set location = "Coffee machine";
   set model = "HP 8860";
   set networkName = "PrinterDEF";
}

computer XYZ {
   set location = "Mike's desk";
   set 'model' = "Mac Book Air";
}
`

func Test_SimpleDeclaration(t *testing.T) {
	RegisterTestingT(t)

	r := strings.NewReader("set x = 5 \"Description\";")
	elems, errs, success := Parse(r)

	Expect(success).To(BeTrue())
	Expect(len(errs)).To(Equal(0))
	Expect(len(elems)).To(Equal(1))

	elem := elems[0]

	Expect(elem.isDeclaration()).To(BeTrue())
	Expect(elem.isDefinition()).To(BeFalse())
	Expect(len(elem.Modifications)).To(Equal(0))

	Expect(elem.Qualifiers).To(Equal([]string{"set"}))
	Expect(elem.Name).To(Equal("x"))
	Expect(elem.Description).To(Equal("Description"))
	Expect(elem.Value).To(Equal(5))
}

func Test_Errors(t *testing.T) {
	RegisterTestingT(t)
}

func Test_SampleInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader("set x = 5")

	_, errs, success := Parse(r)
	Expect(success).To(BeFalse())
	Expect(len(errs)).To(Equal(1))
	Expect(errs[0].Error()).To(Equal("Error syntax error at line 0, column 9"))
}

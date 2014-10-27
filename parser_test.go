package denada

import "testing"
import "strings"

import . "github.com/onsi/gomega"

var sample_noexprs = `
class ABC() "D1" {
   Real foo;
   Integer x;
}

class DEF "D2" {
   String y();
   Boolean x "bool";
}
`

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

func Test_YYSimpleDeclaration(t *testing.T) {
	RegisterTestingT(t)

	r := strings.NewReader("set x = 5 \"Description\";")
	elems, err := Parse(r)

	Expect(err).To(BeNil())
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

func Test_YYErrors(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader("set x = 5")

	exp := "Parsing errors:\n  Error syntax error at line 0, column 9"
	_, err := Parse(r)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal(exp))
}

func Test_YYSampleInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader(sample)

	el, err := Parse(r)
	Expect(err).To(BeNil())

	Expect(len(el)).To(Equal(3))
	Expect(el[0].isDefinition()).To(BeTrue())
	Expect(el[1].isDefinition()).To(BeTrue())
	Expect(el[2].isDefinition()).To(BeTrue())
}

func Test_YYSampleNoExprInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader(sample_noexprs)

	el, err := Parse(r)
	Expect(err).To(BeNil())

	Expect(len(el)).To(Equal(2))
	Expect(el[0].isDefinition()).To(BeTrue())
	Expect(el[1].isDefinition()).To(BeTrue())
}

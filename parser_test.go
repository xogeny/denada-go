package denada

import "os"
import "log"
import "testing"
import "strings"

import . "github.com/onsi/gomega"

var plog = log.New(os.Stderr, "", log.LstdFlags)

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

var sample_exprs = `
Real x = 5.0;
Integer y = 1;
String z = "This is a \"test\"";
Object a = {"key1": 5, "\"test\"": 2, "nested": {"r": "another string"}};
Null b = null;
Boolean c = [true, false];
Array d = [{"x": 5}, "foo", "\"test\"", true];
class Foo(x=5.0, y=1, z="This is a \"test\"", a={"key1": 5}, b=null, c=[true, false],
          d = [{"x": 5}, "foo"]) {
  Real x = 5.0;
}
`

func Test_LLSimpleDeclaration(t *testing.T) {
	RegisterTestingT(t)

	r := strings.NewReader("set x = 5 \"Description\";")
	p, err := NewParser(r, plog)
	elems, err := p.ParseFile()

	Expect(err).To(BeNil())
	Expect(len(elems)).To(Equal(1))

	elem := elems[0]

	Expect(elem.isDeclaration()).To(BeTrue())
	Expect(elem.isDefinition()).To(BeFalse())
	Expect(len(elem.Modifications)).To(Equal(0))

	Expect(elem.Qualifiers).To(Equal([]string{"set"}))
	Expect(elem.Name).To(Equal("x"))
	Expect(elem.Description).To(Equal("Description"))
	v, err := elem.Value.Int()
	Expect(err).To(BeNil())
	Expect(v).To(Equal(5))
}

func Test_LLErrors(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader("set x = 5")

	p, err := NewParser(r, plog)
	_, err = p.ParseFile()

	Expect(err).ToNot(BeNil())
}

func Test_LLSampleInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader(sample)

	p, err := NewParser(r, plog)
	el, err := p.ParseFile()

	Expect(err).To(BeNil())

	Expect(len(el)).To(Equal(3))
	Expect(el[0].isDefinition()).To(BeTrue())
	Expect(el[1].isDefinition()).To(BeTrue())
	Expect(el[2].isDefinition()).To(BeTrue())
}

func Test_LLSampleNoExprInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader(sample_noexprs)

	p, err := NewParser(r, plog)
	el, err := p.ParseFile()

	Expect(err).To(BeNil())

	Expect(len(el)).To(Equal(2))
	Expect(el[0].isDefinition()).To(BeTrue())
	Expect(el[1].isDefinition()).To(BeTrue())
}

func Test_LLSampleJSONInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader(sample_exprs)

	p, err := NewParser(r, plog)
	el, err := p.ParseFile()

	Expect(err).To(BeNil())

	Expect(len(el)).To(Equal(8))
	for i, e := range el {
		if i == 7 {
			Expect(e.isDefinition()).To(BeTrue())
		} else {
			Expect(e.isDefinition()).To(BeFalse())
		}
	}
}

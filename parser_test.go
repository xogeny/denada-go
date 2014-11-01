package denada

import "testing"
import "strings"
import "encoding/json"

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

func Test_SimpleDeclaration(t *testing.T) {
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
	v, err := elem.Value.Int()
	Expect(err).To(BeNil())
	Expect(v).To(Equal(5))
}

func Test_Errors(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader("set x = 5")

	_, err := Parse(r)

	Expect(err).ToNot(BeNil())
}

func Test_SampleInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader(sample)

	el, err := Parse(r)

	Expect(err).To(BeNil())

	Expect(len(el)).To(Equal(3))
	Expect(el[0].isDefinition()).To(BeTrue())
	Expect(el[1].isDefinition()).To(BeTrue())
	Expect(el[2].isDefinition()).To(BeTrue())
}

func Test_SampleNoExprInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader(sample_noexprs)

	el, err := Parse(r)

	Expect(err).To(BeNil())

	Expect(len(el)).To(Equal(2))
	Expect(el[0].isDefinition()).To(BeTrue())
	Expect(el[1].isDefinition()).To(BeTrue())
}

func Test_JsonTypes(t *testing.T) {
	var expr interface{}
	str := `{"minItems": 1}`
	err := json.Unmarshal([]byte(str), &expr)
	Expect(err).To(BeNil())
	asmap, ok := expr.(map[string]interface{})
	Expect(ok).To(Equal(true))
	v, exists := asmap["minItems"]
	Expect(exists).To(Equal(true))
	Expect(v).To(Equal(1.0))
}

func Test_NumbersInExpr(t *testing.T) {
	elems, err := ParseString(`var x = {"minItems": 1 };`)
	Expect(err).To(BeNil())
	e := elems[0]
	v := e.Value

	asmap, err := v.Map()
	Expect(err).To(BeNil())

	mif, exists := asmap["minItems"]
	Expect(exists).To(Equal(true))
	Expect(mif).To(Equal(1.0))

	mi := v.Get("minItems").MustInt()
	Expect(mi).To(Equal(1))
}

func Test_SampleJSONInput(t *testing.T) {
	RegisterTestingT(t)
	r := strings.NewReader(sample_exprs)

	el, err := Parse(r)

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

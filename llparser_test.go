package denada

import "os"
import "log"
import "testing"
import "strings"

import . "github.com/onsi/gomega"

var plog = log.New(os.Stderr, "", log.LstdFlags)

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
	Expect(elem.Value).To(Equal(5))
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

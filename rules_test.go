package denada

import "testing"
import . "github.com/onsi/gomega"

func TestSingularRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("singleton")
	Expect(err).To(BeNil())
	Expect(info.Recursive).To(BeFalse())
	Expect(info.Name).To(Equal("singleton"))
	Expect(info.Cardinality).To(Equal(Cardinality(Singleton)))
}

func TestOptionalRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("optional?")
	Expect(err).To(BeNil())
	Expect(info.Recursive).To(BeFalse())
	Expect(info.Name).To(Equal("optional"))
	Expect(info.Cardinality).To(Equal(Cardinality(Optional)))
}

func TestZoMRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("zom*")
	Expect(err).To(BeNil())
	Expect(info.Recursive).To(BeFalse())
	Expect(info.Name).To(Equal("zom"))
	Expect(info.Cardinality).To(Equal(Cardinality(ZeroOrMore)))
}

func TestOoMRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("oom+")
	Expect(err).To(BeNil())
	Expect(info.Recursive).To(BeFalse())
	Expect(info.Name).To(Equal("oom"))
	Expect(info.Cardinality).To(Equal(Cardinality(OneOrMore)))
}

func TestRecursiveRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("^recur")
	Expect(err).To(BeNil())
	Expect(info.Recursive).To(BeTrue())
	Expect(info.Name).To(Equal("recur"))
	Expect(info.Cardinality).To(Equal(Cardinality(Singleton)))
}

func TestRecursiveComplexRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("^recur?")
	Expect(err).To(BeNil())
	Expect(info.Recursive).To(BeTrue())
	Expect(info.Name).To(Equal("recur"))
	Expect(info.Cardinality).To(Equal(Cardinality(Optional)))
}

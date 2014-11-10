package denada

import "testing"
import . "github.com/onsi/gomega"

func TestSingularRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("singleton", emptyContext)
	Expect(err).To(BeNil())
	Expect(info.Contents).To(Equal(ElementList{}))
	Expect(info.Name).To(Equal("singleton"))
	Expect(info.Cardinality).To(Equal(Cardinality(Singleton)))
}

func TestOptionalRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("optional?", emptyContext)
	Expect(err).To(BeNil())
	Expect(info.Contents).To(Equal(ElementList{}))
	Expect(info.Name).To(Equal("optional"))
	Expect(info.Cardinality).To(Equal(Cardinality(Optional)))
}

func TestZoMRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("zom*", emptyContext)
	Expect(err).To(BeNil())
	Expect(info.Contents).To(Equal(ElementList{}))
	Expect(info.Name).To(Equal("zom"))
	Expect(info.Cardinality).To(Equal(Cardinality(ZeroOrMore)))
}

func TestOoMRule(t *testing.T) {
	RegisterTestingT(t)

	info, err := ParseRule("oom+", emptyContext)
	Expect(err).To(BeNil())
	Expect(info.Contents).To(Equal(ElementList{}))
	Expect(info.Name).To(Equal("oom"))
	Expect(info.Cardinality).To(Equal(Cardinality(OneOrMore)))
}

func TestRecursiveRule(t *testing.T) {
	RegisterTestingT(t)

	root := ElementList{new(Element)}
	context := map[string]ElementList{"$root": root}

	info, err := ParseRule("recur>$root", context)
	Expect(err).To(BeNil())
	Expect(info.Contents).To(Equal(root))
	Expect(info.Name).To(Equal("recur"))
	Expect(info.Cardinality).To(Equal(Cardinality(Singleton)))
}

func TestRecursiveComplexRule(t *testing.T) {
	RegisterTestingT(t)

	root := ElementList{new(Element)}
	context := map[string]ElementList{"$root": root}

	info, err := ParseRule("recur?>$root", context)
	Expect(err).To(BeNil())
	Expect(info.Contents).To(Equal(root))
	Expect(info.Name).To(Equal("recur"))
	Expect(info.Cardinality).To(Equal(Cardinality(Optional)))
}

package denada

import "testing"
import . "github.com/onsi/gomega"

func Test_QualifierMatch(t *testing.T) {
	RegisterTestingT(t)

	g := Element{Qualifiers: []string{"set"}, Name: "_", Description: "foo*", definition: false}
	i := Element{Qualifiers: []string{"var"}, Name: "x", definition: false}

	m := matchQualifiers(&i, &g)
	Expect(m).To(Equal(false))

	gl := ElementList{&g}
	il := ElementList{&i}

	ml := Check(il, gl, false)
	Expect(ml).ToNot(BeNil())
}

func Test_StringMatch(t *testing.T) {
	RegisterTestingT(t)

	match := matchString("abc", "abc")
	Expect(match).To(BeTrue())
	match = matchString("abcabc", "(abc)+")
	Expect(match).To(BeTrue())
	match = matchString("abc", "_")
	Expect(match).To(BeTrue())
	match = matchString("abc", ".+")
	Expect(match).To(BeTrue())

	match = matchString("abc", "def")
	Expect(match).To(BeFalse())
	match = matchString("_", "abc")
	Expect(match).To(BeFalse())
}

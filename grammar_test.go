package denada

import "testing"
import . "github.com/onsi/gomega"

var config_grammar = `
section _ "section*" {
  set _ = _ "variable*";
}`

var config_input1 = `
section Authentication {
  set username = "foo";
  set password = "bar";
}

section DNS {
  set hostname = "localhost";
  set MTU = 1500;
}
`

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

func Test_Grammar(t *testing.T) {
	RegisterTestingT(t)

	il, ies, is := ParseString(config_input1)
	Expect(is).To(BeTrue())
	Expect(ies).To(Equal([]error{}))

	gl, ges, gs := ParseString(config_grammar)
	Expect(gs).To(BeTrue())
	Expect(ges).To(Equal([]error{}))

	errs := Check(il, gl)
	Expect(errs).To(Equal([]error{}))
}

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

var config_err1 = `
section section Authentication {
  set username = "foo";
  set password = "bar";
}

section DNS {
  set hostname = "localhost";
  set MTU = 1500;
}
`

var config_err2 = `
extra section Authentication {
  set username = "foo";
  set password = "bar";
}

section DNS {
  set hostname = "localhost";
  set MTU = 1500;
}
`

var config_err3 = `
section Authentication {
  set username = "foo";
  set password = "bar";
}

section DNS {
  var hostname = "localhost";
  set MTU = 1500;
}
`

func Test_QualifierMatch(t *testing.T) {
	g := Element{Qualifiers: []string{"set"}, Name: "_", Description: "foo*", definition: false}
	i := Element{Qualifiers: []string{"var"}, Name: "x", definition: false}

	m := matchQualifiers(&i, &g)
	Expect(m).To(Equal(false))

	gl := ElementList{&g}
	il := ElementList{&i}

	ml := Check(il, gl)
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

func Test_Grammar(t *testing.T) {
	RegisterTestingT(t)

	gl, ge := ParseString(config_grammar)
	Expect(ge).To(BeNil())

	il, ie := ParseString(config_input1)
	Expect(ie).To(BeNil())

	err := Check(il, gl)
	Expect(err).To(BeNil())

	e1, e1e := ParseString(config_err1)
	Expect(e1e).To(BeNil())

	err = Check(e1, gl)
	Expect(err).ToNot(BeNil())

	e2, e2e := ParseString(config_err2)
	Expect(e2e).To(BeNil())

	err = Check(e2, gl)
	Expect(err).ToNot(BeNil())

	e3, e3e := ParseString(config_err3)
	Expect(e3e).To(BeNil())

	err = Check(e3, gl)
	Expect(err).ToNot(BeNil())
}

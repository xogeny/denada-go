package denada

import "log"
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
	log.Printf("ml: %v", ml)
	Expect(len(ml)).To(Equal(1))
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

	gl, ges, gs := ParseString(config_grammar)
	Expect(gs).To(BeTrue())
	Expect(ges).To(Equal([]error{}))

	il, ies, is := ParseString(config_input1)
	Expect(is).To(BeTrue())
	Expect(ies).To(Equal([]error{}))

	errs := Check(il, gl)
	Expect(errs).To(Equal([]error{}))

	e1, _, e1s := ParseString(config_err1)
	Expect(e1s).To(BeTrue())
	Expect(len(Check(e1, gl))).ToNot(Equal(0))

	e2, _, e2s := ParseString(config_err2)
	Expect(e2s).To(BeTrue())
	Expect(len(Check(e2, gl))).ToNot(Equal(0))

	e3, _, e3s := ParseString(config_err3)
	Expect(e3s).To(BeTrue())
	Expect(len(Check(e3, gl))).ToNot(Equal(0))
}

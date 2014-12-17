package denada

import "fmt"
import "testing"
import . "github.com/onsi/gomega"

type Case1 struct {
	//X     string
	Reals []struct {
		Named string `dnd:"name=_"`
		Label string `dnd:"mod"`
		Units string `dnd:"mod"`
		Value int    `dnd:"value"`
	} `dndrule:"Real*",dndquals:"Real"`
	Groups []struct {
		Named    string  `dnd:"name"`
		Label    string  `dnd:"mod"`
		Image    string  `dnd:"mod"`
		Contents []Case1 `dnd:"contents"`
	} `dndrule:"Groups*" dndquals:"group"`
}

type DenadaFormat interface {
	Grammar() ElementList
	Unmarshal(elems ElementList) error
}

/*
type NamedStringVar struct {
	Name  string
	Value string
}

func (n NamedStringVar) Grammar() Element {
	return Element{
		Name:  n.Name,
		Value: "$string",
	}
}

type Case2 struct {
	project   NamedStringVar
	processes []NamedStringVar
}
*/

func Test_MarshalCase1(t *testing.T) {
	RegisterTestingT(t)

	c := Case1{}
	elems, err := Marshal(c)
	fmt.Printf("%s\n", Unparse(elems))
	Expect(err).To(BeNil())
	Expect(len(elems)).To(Equal(2))
}

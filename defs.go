package denada

import "fmt"

import "github.com/bitly/go-simplejson"

type Modifications map[string]*simplejson.Json

type Element struct {
	/* Common to all elements */
	Qualifiers    []string
	Name          string
	Description   string
	Modifications Modifications

	/* For definitions */
	Contents ElementList

	/* For declarations */
	Value *simplejson.Json

	rule       string
	definition bool
}

func (e Element) String() string {
	ret := ""
	for _, q := range e.Qualifiers {
		ret += q + " "
	}
	ret += e.Name

	if e.isDefinition() {
		return fmt.Sprintf("%s { ... }", ret)
	} else {
		if e.Value != nil {
			return fmt.Sprintf("%s = %v;", ret, e.Value)
		} else {
			return fmt.Sprintf("%s;", ret)
		}
	}
}

func (e Element) Rule() string {
	return e.rule
}

func (e Element) isDefinition() bool {
	return e.definition
}

func (e Element) isDeclaration() bool {
	return !e.definition
}

type ElementList []*Element

func (e ElementList) Definition(name string) (*Element, error) {
	for _, d := range e {
		if d.isDefinition() && d.Name == name {
			return d, nil
		}
	}
	return nil, fmt.Errorf("Unable to find definition for %s", name)
}

func MakeElementList() ElementList {
	return ElementList{}
}

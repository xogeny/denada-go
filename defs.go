package denada

import "fmt"

type Expr interface{}

type Element struct {
	/* Common to all elements */
	Qualifiers    []string
	Name          string
	Description   string
	Modifications map[string]Expr

	/* For definitions */
	Contents ElementList

	/* For declarations */
	Value Expr

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
	return nil, fmt.Errorf("Unable to find definition for %s")
}

func MakeElementList() ElementList {
	return ElementList{}
}

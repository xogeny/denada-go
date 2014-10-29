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
	Contents      ElementList      // Used by definitions
	Value         *simplejson.Json // Used by declarations

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

func (e ElementList) Definitions() ElementList {
	ret := ElementList{}
	for _, elem := range e {
		if elem.isDefinition() {
			ret = append(ret, elem)
		}
	}
	return ret
}

func (e ElementList) Declarations() ElementList {
	ret := ElementList{}
	for _, elem := range e {
		if elem.isDeclaration() {
			ret = append(ret, elem)
		}
	}
	return ret
}

func (e ElementList) Named(name string) ElementList {
	ret := ElementList{}
	for _, elem := range e {
		if elem.Name == name {
			ret = append(ret, elem)
		}
	}
	return ret
}

func (e ElementList) AllElements() ElementList {
	ret := ElementList{}
	for _, elem := range e {
		ret = append(ret, elem)
		if elem.isDefinition() {
			ret = append(ret, elem.Contents.AllElements()...)
		}
	}
	return ret
}

func (e ElementList) PopHead() (*Element, ElementList, error) {
	if len(e) == 0 {
		return nil, e, fmt.Errorf("Cannot pop the head of an empty element list")
	}
	ret := e[0]
	e = e[1:]
	return ret, e, nil
}

func MakeElementList() ElementList {
	return ElementList{}
}

package denada

import (
	"fmt"
	"github.com/bitly/go-simplejson"
)

type ElementList []*Element

func (e ElementList) Definition(name string, children ...string) (*Element, error) {
	for _, d := range e {
		if d.IsDefinition() && d.Name == name {
			if len(children) == 0 {
				return d, nil
			} else {
				return d.Contents.Definition(children[0], children[1:]...)
			}
		}
	}
	return nil, fmt.Errorf("Unable to find definition for %s", name)
}

func (e ElementList) Definitions() ElementList {
	ret := ElementList{}
	for _, elem := range e {
		if elem.IsDefinition() {
			ret = append(ret, elem)
		}
	}
	return ret
}

func (e ElementList) Declarations() ElementList {
	ret := ElementList{}
	for _, elem := range e {
		if elem.IsDeclaration() {
			ret = append(ret, elem)
		}
	}
	return ret
}

func (e ElementList) FirstNamed(name string) *Element {
	for _, elem := range e {
		if elem.Name == name {
			return elem
		}
	}
	return nil
}

func (e ElementList) QualifiedWith(name ...string) ElementList {
	ret := ElementList{}
	for _, elem := range e {
		if elem.HasQualifiers(name...) {
			ret = append(ret, elem)
		}
	}
	return ret
}

func (e ElementList) OfRule(name string, fqn bool) ElementList {
	ret := ElementList{}
	for _, elem := range e {
		if fqn {
			if elem.rulepath == name {
				ret = append(ret, elem)
			}
		} else {
			if elem.rule == name {
				ret = append(ret, elem)
			}
		}
	}
	return ret
}

func (e ElementList) AllElements() ElementList {
	ret := ElementList{}
	for _, elem := range e {
		ret = append(ret, elem)
		if elem.IsDefinition() {
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

func (e ElementList) Equals(o ElementList) error {
	// Now make sure they have the same number of children
	if len(e) != len(o) {
		return fmt.Errorf("Mismatch in number of child elements: %d vs %d: %v vs %v",
			len(e), len(o), e, o)
	}

	// And that each child is equal
	for cn, child := range e {
		err := child.Equals(*o[cn])
		if err != nil {
			return err
		}
	}
	return nil
}

// GetValue tries to find an element that matches the **fully qualified** rule name
// provided.  It returns nil if it cannot find a match (or it finds multiple matches),
// otherwise it returns the value associated with that declaration.
func (e ElementList) GetValue(rulename string) *simplejson.Json {
	elems := e.AllElements().OfRule(rulename, true)
	if len(elems) != 1 {
		return nil
	}
	return elems[0].Value
}

func (e ElementList) GetStringValue(rulename string, defaultValue string) string {
	v := e.GetValue(rulename)
	if v == nil {
		return defaultValue
	}
	ret := v.MustString()
	if ret == "" {
		return defaultValue
	}
	return ret
}

func (e ElementList) GetIntValue(rulename string, defaultValue int) int {
	v := e.GetValue(rulename)
	if v == nil {
		return defaultValue
	}
	ret, err := v.Int()
	if err != nil {
		return defaultValue
	}
	return ret
}

func (e ElementList) GetBoolValue(rulename string, defaultValue bool) bool {
	v := e.GetValue(rulename)
	if v == nil {
		return defaultValue
	}
	ret, err := v.Bool()
	if err != nil {
		return defaultValue
	}
	return ret
}

func MakeElementList() ElementList {
	return ElementList{}
}

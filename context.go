package denada

import "log"
import "fmt"

type RuleContext struct {
	This   ElementList
	parent *RuleContext
}

func NullContext() RuleContext {
	return RuleContext{
		This:   ElementList{},
		parent: nil,
	}
}

func RootContext(elems ElementList) RuleContext {
	return RuleContext{
		This:   elems,
		parent: nil,
	}
}

func ChildContext(elems ElementList, parent *RuleContext) RuleContext {
	return RuleContext{
		This:   elems,
		parent: parent,
	}
}

func (c RuleContext) String() string {
	names := []string{}
	for _, e := range c.This {
		rule, err := ParseRuleName(e.Description)
		if err != nil {
			names = append(names, "<non-rule>")
		} else {
			names = append(names, rule.Name)
		}
	}
	if c.parent == nil {
		return fmt.Sprintf("[%s]", names)
	} else {
		return fmt.Sprintf("[%s (%v)]", names, c.parent)
	}
}

func (c RuleContext) Find(path ...string) (ret RuleContext, err error) {
	/* If no path is provided, they must me this context */
	if len(path) == 0 {
		ret = c
		return
	}

	/* Take the first element of the path and check against some "reserved" names */
	head := path[0]

	/* Are they looking in the current context? */
	if head == "." {
		return c.Find(path[1:]...)
	}

	/* Are they looking for the parent context? */
	if head == ".." {
		if c.parent == nil {
			/* If we are at the root, we have no parent */
			err = fmt.Errorf("Requested rule context at root level")
			return
		} else {
			/* Otherwise, resume the search in our parent context */
			log.Printf("Looking up path %v in parent", path)
			return c.parent.Find(path[1:]...)
		}
	}

	/* Are they looking for the root context? */
	if head == "$root" {
		if c.parent == nil {
			/* If we are at the root, search this context */
			return c.Find(path[1:]...)
		} else {
			/* If not, ask our parent for $root */
			return c.parent.Find(path...)
		}
	}

	/* Check to see if head matches any uniquely named child definitions */
	result := ElementList{}
	for _, d := range c.This {
		if !d.IsDefinition() {
			continue
		}
		rule, err := ParseRuleName(d.Description)
		if err != nil {
			continue
		}
		if rule.Name == head {
			result = append(result, d.Contents...)
		}
	}
	if result != nil {
		return ChildContext(result, &c), nil
	}

	err = fmt.Errorf("Unable to find context %s", head)
	return
}

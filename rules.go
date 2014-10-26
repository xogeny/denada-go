package denada

import "fmt"

type Cardinality int

const (
	Zero = iota
	Optional
	ZeroOrMore
	Singleton
	OneOrMore
)

type RuleInfo struct {
	Recursive   bool
	Name        string
	Cardinality Cardinality
}

func (r RuleInfo) checkCount(count int) error {
	switch r.Cardinality {
	case Zero:
		if count != 0 {
			return fmt.Errorf("Expected zero of rule %s, found %d", r.Name, count)
		}
	case Optional:
		if count > 1 {
			return fmt.Errorf("Expected at most 1 of rule %s, found %d", r.Name, count)
		}
	case Singleton:
		if count != 1 {
			return fmt.Errorf("Expected at exactly 1 of rule %s, found %d", r.Name, count)
		}
	case OneOrMore:
		if count == 0 {
			return fmt.Errorf("Expected at least 1 of rule %s, found %d", r.Name, count)
		}
	}
	return nil
}

func ParseRule(e Element) (rule RuleInfo, err error) {
	rule = RuleInfo{Cardinality: Zero}
	if e.Description == "" {
		err = fmt.Errorf("Rule element %s doesn't include a description", e.String())
		return
	}
	str := e.Description
	if str[0] == '^' {
		str = str[1:]
		rule.Recursive = true
	}
	l := len(str) - 1
	lastchar := str[l]
	if lastchar == '+' {
		rule.Cardinality = OneOrMore
		str = str[0:l]
	} else if lastchar == '*' {
		rule.Cardinality = ZeroOrMore
		str = str[0:l]
	} else if lastchar == '?' {
		rule.Cardinality = Optional
		str = str[0:l]
	} else {
		rule.Cardinality = Singleton
	}
	rule.Name = str
	return
}

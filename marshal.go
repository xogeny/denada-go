package denada

import "log"
import "fmt"
import "strings"
import "reflect"

func Marshal(data interface{}) (ElementList, error) {
	ret := []*Element{}
	typ := reflect.TypeOf(data)
	log.Printf("Type: %v", typ)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		log.Printf("  Field %d: %v", i, f)
		e, err := marshalField(f)
		if err != nil {
			err = fmt.Errorf("Error marshalling field #%d: %v", i, err)
			return nil, err
		}
		ret = append(ret, e)
	}

	return ret, nil
}

func marshalField(f reflect.StructField) (elem *Element, err error) {
	rule := f.Tag.Get("dndrule")
	if rule == "" {
		return nil, fmt.Errorf("Field %s has no dndrule field", f.Name)
	}
	_, err = ParseRuleName(rule)
	if err != nil {
		return nil, fmt.Errorf("Invalid rule '%s' associated with field %s",
			rule, f.Name)
	}

	quals := f.Tag.Get("dndquals")
	log.Printf("quals = %s", quals)

	return &Element{
		Qualifiers:  strings.Split(quals, " "),
		Name:        "_",
		Description: rule,
	}, nil
}

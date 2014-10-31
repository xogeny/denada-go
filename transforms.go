package denada

import "fmt"

var importGrammar = `
import(file="$string", recursive?="$bool") "import*";
`

func ImportTransform(root ElementList) (ret ElementList, err error) {
	g, err := ParseString(importGrammar)
	if err != nil {
		err = fmt.Errorf("Error parsing import statement grammar: %v", err)
		return
	}

	ret = ElementList{}
	for _, e := range root {
		match := matchElement(e, g[0], g[0].Contents, false, "")
		if match == nil {
			file := e.Modifications["file"].MustString()
			insert, err := ParseFile(file)
			if err != nil {
				return nil, fmt.Errorf("Error parsing import contents from %s: %v", file, err)
			}
			rval, present := e.Modifications["recursive"]
			recursive := false
			if present {
				recursive = rval.MustBool()
			}

			if recursive {
				insert, err = ImportTransform(insert)
				if err != nil {
					return nil, err
				}
			}
			for _, i := range insert {
				ret = append(ret, i)
			}
		} else {
			if e.isDefinition() {
				newchildren, err := ImportTransform(e.Contents)
				if err != nil {
					return nil, err
				}
				newe := e.Clone()
				newe.Contents = newchildren
				ret = append(ret, newe)
			} else {
				ret = append(ret, e)
			}
		}
	}
	return ret, nil
}

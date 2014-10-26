package denada

type Element struct {
	/* Common to all elements */
	Qualifiers    []string
	Name          string
	Description   string
	Modifications map[string]interface{}

	/* For definitions */
	Contents []Element

	/* For declarations */
	Value interface{}

	definition bool
}

func (e Element) isDefinition() bool {
	return e.definition
}

func (e Element) isDeclaration() bool {
	return !e.definition
}

package denada

import "fmt"
import "regexp"
import "log"

import "github.com/bitly/go-simplejson"
import "github.com/xeipuuv/gojsonschema"

func Check(input ElementList, grammar ElementList, diag bool) error {
	return CheckContents(input, grammar, diag, "", "")
}

type matchInfo struct {
	count int
	rule  RuleInfo
	desc  string
}

func CheckContents(input ElementList, grammar ElementList, diag bool,
	prefix string, parentRule string) error {

	// Initialize data associated with rule matching
	counts := map[string]*matchInfo{}

	// Loop over grammar rules and record counts information
	for _, g := range grammar {
		// Make sure grammar element has a (rule) description
		if g.Description == "" {
			return fmt.Errorf("Grammar element %s has no description", g.String())
		}

		// Parse the rule information from the description
		rule, err := ParseRule(g.Description)

		// If there is an error in the rule description, add an error and
		// skip this grammar element
		if err != nil {
			return fmt.Errorf("Error in rule description: %v", err)
		}

		mi, exists := counts[rule.Name]
		if exists {
			if rule.Name != mi.rule.Name || rule.Recursive != mi.rule.Recursive ||
				rule.Cardinality != mi.rule.Cardinality {
				return fmt.Errorf("Unmatching rules with same name: %s vs %s",
					g.Description, mi.desc)
			}
		} else {
			counts[rule.Name] = &matchInfo{count: 0, rule: rule, desc: g.Description}
		}
	}

	// Now, loop over all the actual input elements and see if they match
	// any of the rules
	for _, in := range input {
		var likely error = nil
		ierrs := []error{}
		for _, g := range grammar {
			// Parse the rule information from the description (ignore error
			// because we already checked that)
			rule, _ := ParseRule(g.Description)

			// Normally, if we find a match in matchElement, we'll use that
			// matching rules contents as the rules to apply to the matching
			// inputs children...
			context := g.Contents

			if rule.Recursive {
				// ...but if the rule is recursive, we choose the same rules
				// as we are currently using at this level
				context = grammar
			}
			path := parentRule + "." + rule.Name
			if parentRule == "" {
				path = rule.Name
			}
			ematch := matchElement(in, g, context, diag, prefix, path)
			if ematch == nil {
				// A match was found, so increment the count for this particular
				// grammar rule
				counts[rule.Name].count++

				// Then check to see if this input has matched any previous rules
				if in.rule == "" {
					// If not, indicate what rule this input matched
					if diag {
						log.Printf("%sInput %s matched %s", prefix, in.String(), rule.Name)
					}
					in.rule = path
				} else {
					if diag {
						log.Printf("%sInput %s already matched %s", prefix, in.String(), in.rule)
					}
					// TODO: Just match first one and move on?
					return fmt.Errorf("Element %s matched rule %s and %s",
						in.String(), in.rule, rule.Name)
				}
			} else {
				if diag {
					log.Printf("%sInput %s did not match %s because\n%s", prefix, in.String(),
						rule.Name, ematch.Error())
				}
				if len(grammar) == 1 {
					return ematch
				}
				if len(in.Qualifiers) == 1 && len(g.Qualifiers) == 1 &&
					in.Qualifiers[0] == g.Qualifiers[0] {
					likely = ematch
				}
				if in.Name == g.Name {
					likely = ematch
				}
				ierrs = append(ierrs, ematch)
			}
		}
		if in.rule == "" {
			if likely == nil {
				return fmt.Errorf("No match for element %v because %v", in, listToError(ierrs))
			} else {
				return likely
			}
		}
	}

	// Check to make sure that all rules were matched the correct number
	// of times.
	for _, mi := range counts {
		rerrs := []error{}
		err := mi.rule.checkCount(mi.count)
		if err != nil {
			rerrs = append(rerrs, err)
		}
		if len(rerrs) > 0 {
			return listToError(rerrs)
		}
	}

	return nil
}

func matchString(input string, grammar string) bool {
	if grammar == "_" {
		return true
	}
	matched, err := regexp.MatchString(grammar, input)
	if err == nil && matched {
		return true
	}
	return false
}

func matchQualifiers(input *Element, grammar *Element) bool {
	imatch := make([]bool, len(input.Qualifiers))
	for _, g := range grammar.Qualifiers {
		count := 0

		rule, err := ParseRule(g)

		if err != nil {
			log.Printf("Error parsing rule information in qualifier '%s': %v", g, err)
			return false
		}

		for i, in := range input.Qualifiers {
			matched := matchString(in, rule.Name)
			if matched {
				imatch[i] = true
				count++
			}
		}

		// Check to see if the correct number of matches were found for this qualifier
		err = rule.checkCount(count)
		if err != nil {
			// If not, this is not a match
			return false
		}
	}

	// Now check to make sure every qualifier on the input element had a match
	for i, _ := range input.Qualifiers {
		if !imatch[i] {
			// This qualifier on the input element was never matched
			return false
		}
	}

	return true
}

func matchModifications(input *Element, grammar *Element, diag bool) bool {
	// Create a map to keep track of which modification keys on the input
	// element find a match
	imatch := map[string]bool{}
	for k, _ := range input.Modifications {
		imatch[k] = false
	}

	// Now loop over all keys and expresions in the grammar
	for r, ge := range grammar.Modifications {
		count := 0

		// Parse the rule
		rule, err := ParseRule(r)

		if err != nil {
			// If the rule is not valid, assume no match
			log.Printf("Error parsing rule information in key '%s': %v", r, err)
			return false
		}

		// Loop over all actual modification keys and values
		for i, ie := range input.Modifications {
			// Check to see if the keys match
			matched := matchString(i, rule.Name)
			if matched {
				// If so, check if the expressions match
				if matchExpr(ie, ge, diag) {
					// If so, this input is matched and so is the grammar rule
					imatch[i] = true
					count++
				}
			}
		}

		// Now check to make sure this grammar rule has been matched an appropriate
		// number of times
		err = rule.checkCount(count)
		if err != nil {
			// If not, no match
			return false
		}
	}

	// Now check to make sure every key on the input element had a match
	for k, _ := range input.Modifications {
		if !imatch[k] {
			// This key on the input element was never matched
			return false
		}
	}

	return true
}

// Validation rules:
//   Grammar expr is:
//     String that starts with $ -> Look for type match
//     String (without $) -> Exact match
//     Object -> Treat object as a JSON schema and validate input with it
//     Otherwise -> No match
func matchExpr(input *simplejson.Json, grammar *simplejson.Json, diag bool) bool {
	if grammar == nil && input == nil {
		return true
	}
	if grammar == nil || input == nil {
		if diag {
			log.Printf("Grammar was %v while input was %v", grammar, nil)
		}
		return false
	}
	stype, err := grammar.String()
	if err == nil {
		switch stype {
		case "$_":
			return true
		case "$string":
			_, terr := input.String()
			if terr != nil && diag {
				log.Printf("Input wasn't a string")
			}
			return terr == nil
		case "$bool":
			_, terr := input.Bool()
			if terr != nil && diag {
				log.Printf("Input wasn't a bool")
			}
			return terr == nil
		case "$int":
			_, terr := input.Int64()
			if terr != nil && diag {
				log.Printf("Input wasn't an int")
			}
			return terr == nil
		case "$number":
			_, terr := input.Float64()
			if terr != nil && diag {
				log.Printf("Input wasn't a number")
			}
			return terr == nil
		default:
			is, terr := input.String()
			log.Printf("treated as literal")
			return terr == nil && is == stype
		}
	}
	mtype, err := grammar.Map()
	if err == nil {
		schema, err := gojsonschema.NewJsonSchemaDocument(mtype)
		if err != nil {
			log.Printf("Invalid schema (%v) in grammar: %v", mtype, err)

			return false
		}
		result := schema.Validate(input)

		for _, e := range result.Errors() {
			log.Printf("  JSON Schema validation failed because: %s", e)
		}
		return result.Valid()
	}
	return false
}

func matchElement(input *Element, grammar *Element,
	context ElementList, diag bool, prefix string, parentRule string) error {
	// Check if the names match
	matched := matchString(input.Name, grammar.Name)

	// If the names don't match, no match
	if !matched {
		return fmt.Errorf("Name mismatch (%s doesn't match pattern %s)",
			input.Name, grammar.Name)
	}

	// Check whether the input is a definition or declaration
	if input.isDefinition() {
		if grammar.isDeclaration() {
			// If the input is a definition but the grammar is a declaration, no match
			return fmt.Errorf("Element type mismatch between %v and %v", input, grammar)
		}
		cerr := CheckContents(input.Contents, context, diag, prefix+"  ", parentRule)
		if cerr != nil {
			// If the contents of input don't match the contents of grammar, no match
			return cerr
		}
	} else {
		if grammar.isDefinition() {
			// If the input is a declaration but the grammar is a definition, no match
			return fmt.Errorf("Element type mismatch")
		}
		if !matchExpr(input.Value, grammar.Value, diag) {
			return fmt.Errorf("Value pattern mismatch")
		}
	}

	if !matchQualifiers(input, grammar) {
		return fmt.Errorf("Qualifier mismatch (%v vs %v)", input.Qualifiers,
			grammar.Qualifiers)
	}

	if !matchModifications(input, grammar, diag) {
		return fmt.Errorf("Modification mismatch (%v vs %v)", input.Modifications,
			grammar.Modifications)
	}

	return nil
}

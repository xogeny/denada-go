package denada

import "fmt"

func Check(input ElementList, grammar ElementList) []error {
	// Create a list of errors for this context
	ret := []error{}

	// Loop over grammar rules
	for _, g := range grammar {
		// Initialize how many matches have been made for this rule
		count := 0

		// Make sure grammar element has a (rule) description
		if g.Description == "" {
			ret = append(ret, fmt.Errorf("Grammar element %s has no description", g.String()))
			continue
		}

		// Parse the rule information from the description
		rule, err := ParseRule(g.Description)

		// If there is an error in the rule description, add an error and
		// skip this grammar element
		if err != nil {
			ret = append(ret, err)
			continue
		}

		// Now, loop over all the actual input elements and see if they match
		// any of the rules
		for _, in := range input {
			if match(in, g) {
				// A match was found, so increment the count for this particular
				// grammar rule
				count++

				// Then check to see if this input has matched any previous rules
				if in.rule == "" {
					// If not, indicate what rule this input matched
					in.rule = rule.Name
				} else {
					// If so, add an error
					ret = append(ret, fmt.Errorf("Element %s matched rule %s and %s",
						in.String(), in.rule, rule.Name))
				}
			}
		}

		// Now that we have checked all inputs in this context, check to see if the
		// cardinality of this grammar element was met
		err = rule.checkCount(count)
		if err != nil {
			// If not, add an error
			ret = append(ret, err)
		}
	}

	// Finally, look over all input elements and make sure they matched at least
	// one grammar rule
	for _, in := range input {
		if in.rule == "" {
			// If not, add an error
			ret = append(ret, fmt.Errorf("Element %s didn't match any rule", in.String()))
		}
	}

	// Return any errors that were found in this context
	return ret
}

func match(input Element, grammar Element) bool {
	return false
}

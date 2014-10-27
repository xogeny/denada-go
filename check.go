package denada

import "fmt"
import "regexp"
import "log"

func Check(input ElementList, grammar ElementList) []error {
	// Create a list of errors for this context
	ret := []error{}

	// TODO: Handle multiple grammar nodes with the same rule name

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

func matchQualifiers(input Element, grammar Element) bool {
	imatch := make([]bool, len(input.Qualifiers))
	for _, g := range grammar.Qualifiers {
		count := 0

		rule, err := ParseRule(g)

		if err != nil {
			log.Printf("Error parsing rule information in qualifier '%s': %v", g, err)
			return false
		}

		for i, in := range input.Qualifiers {
			matched, err := regexp.MatchString(rule.Name, in)
			if err == nil && matched {
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

func matchModifications(input Element, grammar Element) bool {
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
			matched, err := regexp.MatchString(rule.Name, i)
			if err == nil && matched {
				// If so, check if the expressions match
				if matchExpr(ie, ge) {
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

// TODO: Implement this (use JSON schema?)
func matchExpr(input Expr, grammar Expr) bool {
	return true
}

func match(input Element, grammar Element) bool {
	// Check if the names match
	matched, err := regexp.MatchString(grammar.Name, input.Name)

	// If the names don't match, no match
	if err != nil || !matched {
		return false
	}

	// Check whether the input is a definition or declaration
	if input.isDefinition() {
		if grammar.isDeclaration() {
			// If the input is a definition but the grammar is a declaration, no match
			return false
		}
		cerrors := Check(input.Contents, grammar.Contents)
		if len(cerrors) > 0 {
			// If the contents of input don't match the contents of grammar, no match
			return false
		}
	} else {
		if grammar.isDefinition() {
			// If the input is a declaration but the grammar is a definition, no match
			return false
		}
		if !matchExpr(input.Value, grammar.Value) {
			return false
		}
	}

	if !matchQualifiers(input, grammar) {
		return false
	}

	if !matchModifications(input, grammar) {
		return false
	}

	return true
}

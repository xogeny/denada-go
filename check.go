package denada

func Check(input ElementList, grammar ElementList) []error {
	ret := []error{}
	for _, g := range grammar {
		count := 0
		rule, err := ParseRule(g)
		if err != nil {
			ret = append(ret, err)
			continue
		}

		for _, in := range input {
			if match(in, g) {
				count++
				break
			}
		}

		err = rule.checkCount(count)
		if err != nil {
			ret = append(ret, err)
		}
	}
	return ret
}

func match(input Element, grammar Element) bool {
	return false
}

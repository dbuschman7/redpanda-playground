package parser

func MultilineParser(snipe rune, predicate Parser[bool]) Parser[[]string] {

	p := Parser[[]string](func(state State) ([]string, State, error) {
		begins := []State{}

		for _, snipe := range state.Tokenize(true, func(r rune) bool { return r == snipe }) {

			//fmt.Printf("Found %v - %v \n", snipe.start, snipe.remaining())

			result, _, err := predicate(snipe)
			if err == nil && result {
				begins = append(begins, snipe)
			} else {
				last := begins[len(begins)-1]
				begins = begins[:len(begins)-1]
				begins = append(begins, last.Extend(snipe))
			}

		}

		if len(begins) == 0 {
			return []string{}, state.consume(state.end), nil
		} else {
			lines := []string{}
			for _, begin := range begins {
				lines = append(lines, begin.remaining())
			}
			return lines, begins[len(begins)-1], nil
		}
	})

	return p
}

package parser

import (
	"fmt"
	"slices"
)

var quotesList = []rune{'"', '\'', '“', '”', '‘', '’', '`'}

func QuotedStringParser() Parser[string] {
	return func(initial State) (string, State, error) {

		quote, start := initial.nextHeadRuneIf(func(r rune) bool {
			return slices.Contains(quotesList, r)
		})

		current := start
		found := false
		for current.end > current.start && !found {
			char, pos := current.nextHeadRune()
			if char == quote {
				found = true
				return start.extractString(current), pos, nil
			}
			current = pos
		}

		if !found {
			return "", initial, fmt.Errorf("no closing quote")
		}
		return "", initial, nil

	}
}

func TailQuotedStringParser() Parser[string] {
	return func(initial State) (string, State, error) {
		quote, current := initial.nextTailRuneIf(func(r rune) bool {
			return slices.Contains(quotesList, r)
		})

		found := false
		// fmt.Printf("start: %v -> %v\n", current.start, current.end)
		for current.end > current.start && !found {
			char, next := current.nextTailRune()
			// fmt.Printf("char: '%c'  %v -> %v \n", char, next.start, next.end)

			if char == quote {
				found = true
				return initial.tailConsume(1).extractTailString(current), next, nil
			}
			current = next
		}

		if !found {
			return "", initial, fmt.Errorf("no closing quote")
		}
		return "", initial, nil

	}
}

package parser

import (
	"unicode/utf8"
)

// state is the internal representation of parsing state.
type State struct {
	data  string // The input string
	start int    // The current parsing offset into the input string.
	end   int    // The end of the string under consideration
}

func WithState(data string) State {
	return State{data: data, start: 0, end: len(data)}
}

// remaining returns the a string which is just the unconsumed input
func (s State) remaining() string {
	return s.data[s.start:s.end]
}

// consume returns a new state in which the offset pointer is advanced
// by n bytes
func (s State) consume(n int) State {
	if n > 0 {
		if s.start+n > s.end {
			n = s.end - s.start
		}
		s.start += n
	}
	return s
}

func (s State) tailConsume(n int) State {
	if n > 0 {
		if s.end-n < s.start {
			n = s.end - s.start
		}
		s.end -= n
	}
	return s
}

func (s State) nextHeadRuneIf(predicate func(rune) bool) (rune, State) {
	r, w := utf8.DecodeRuneInString(s.remaining())
	if predicate(r) {
		return r, s.consume(w)
	}
	return utf8.RuneError, s
}

// nextHeadRune returns the next rune in the input, as well as a new
// state in which the rune has been consumed.
func (s State) nextHeadRune() (rune, State) {
	r, w := utf8.DecodeRuneInString(s.remaining())
	return r, s.consume(w)
}

func (s State) nextTailRuneIf(predicate func(rune) bool) (rune, State) {
	r, w := utf8.DecodeLastRuneInString(s.remaining())
	if predicate(r) {
		return r, s.tailConsume(w)
	}
	return utf8.RuneError, s
}

// nextTailRune returns the last rune in the input, as well as a new
// state in which the rune has been consumed.
func (s State) nextTailRune() (rune, State) {
	r, w := utf8.DecodeLastRuneInString(s.remaining())
	return r, s.tailConsume(w)
}

func (s State) extractString(next State) string {
	return s.data[s.start:next.start]
}

func (s State) extractTailString(next State) string {
	return s.data[next.end:s.end]
}

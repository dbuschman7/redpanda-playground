package parser

import (
	"fmt"
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

func (s State) length() int {
	return s.end - s.start
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

type StateList []State

// snipe returns a list of states in which the next rune is equal to the
// given rune.
func (s State) Snipe(char rune) StateList {
	var states []State
	for i := s.start; i < s.end; i++ {
		//fmt.Printf("i: %d current: %v  looking for %v \n", i, s.data[i], char)
		if s.data[i] == byte(char) {
			states = append(states, State{data: s.data, start: i, end: s.end})
		}
	}
	return states
}

// tokenize returns a list of states which begin to the
// predicate function.
func (s State) Tokenize(keepDelim bool, predicate func(rune) bool) StateList {
	var indexes []int

	indexes = append(indexes, s.start)

	lastSeen := s.start
	for i := s.start; i < s.end; {
		test := predicate(rune(s.data[i]))
		fmt.Printf("i:%d   current:%v   test:%v    lastSeen:%v \n", i, s.data[i], test, lastSeen)
		if test && i > s.start {
			indexes = append(indexes, i)
		}
		i += 1
		lastSeen = i
	}
	indexes = append(indexes, s.end)
	fmt.Printf("indexes: %v \n", indexes)

	var states []State
	for i := range len(indexes) - 1 {
		state := State{data: s.data, start: indexes[i], end: indexes[i+1]}
		if !keepDelim {
			state = state.trimLeadingPredicate(predicate)
		}
		fmt.Printf("state %v: '%v' \n", i, state.remaining())
		if state.length() > 0 {
			states = append(states, state)
		}

	}
	return states
}

func (s State) trimLeadingPredicate(predicate func(rune) bool) State {
	for i := s.start; i < s.end; i++ {
		if !predicate(rune(s.data[i])) {
			return State{data: s.data, start: i, end: s.end}
		}
	}
	return State{data: s.data, start: s.end, end: s.end}
}

func (s State) ContainsRune(r rune) bool {
	for i := s.start; i < s.end; i++ {
		if rune(s.data[i]) == r {
			return true
		}
	}
	return false
}

func (s State) ContainsAnyRunes(runes []rune) bool {
	for i := s.start; i < s.end; i++ {
		for _, r := range runes {
			if rune(s.data[i]) == r {
				return true
			}
		}
	}
	return false
}

func (s State) Extend(in State) State {
	return State{data: s.data, start: s.start, end: in.end}
}

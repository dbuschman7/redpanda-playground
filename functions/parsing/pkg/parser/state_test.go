package parser

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	var state State
	state = WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	state = state.consume(2)
	assert.Equal(t, "llo", state.remaining())
	r, newState := state.nextHeadRune()
	assert.Equal(t, 'l', r)
	assert.Equal(t, "lo", newState.remaining())

	r, newState = newState.nextHeadRune()
	assert.Equal(t, 'l', r)
	assert.Equal(t, "o", newState.remaining())

	r, newState = newState.nextHeadRune()
	assert.Equal(t, 'o', r)
	assert.Equal(t, "", newState.remaining())

	r, newState = newState.nextHeadRune()
	assert.Equal(t, utf8.RuneError, r)
	assert.Equal(t, "", newState.remaining())

	state = WithState("hello world")
	state = state.consume(2)
	assert.Equal(t, "llo world", state.remaining())

	state = state.tailConsume(2)
	assert.Equal(t, "llo wor", state.remaining())

	r, newState = state.nextTailRune()
	assert.Equal(t, 'r', r)
	assert.Equal(t, "llo wo", newState.remaining())

	r, newState = newState.nextTailRuneIf(func(r rune) bool {
		return r == 'o'
	})
	assert.Equal(t, 'o', r)
	assert.Equal(t, "llo w", newState.remaining())

	r, newState = newState.nextTailRuneIf(func(r rune) bool {
		return r == 'q'
	})
	assert.Equal(t, utf8.RuneError, r)
	assert.Equal(t, "llo w", newState.remaining())

	r, newState = newState.nextHeadRuneIf(func(r rune) bool {
		return r == 'l'
	})

	assert.Equal(t, 'l', r)
	assert.Equal(t, "lo w", newState.remaining())

	r, newState = newState.nextHeadRuneIf(func(r rune) bool {
		return r == 'q'
	})
	assert.Equal(t, utf8.RuneError, r)
	assert.Equal(t, "lo w", newState.remaining())

}

func TestStateConsume(t *testing.T) {
	var state State
	state = WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	state = state.consume(2)
	assert.Equal(t, "llo", state.remaining())
}

func TestStateNextHeadRune(t *testing.T) {
	state := WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	r, newState := state.nextHeadRune()
	assert.Equal(t, 'h', r)
	assert.Equal(t, "ello", newState.remaining())
}

func TestStateNextTailRune(t *testing.T) {
	state := WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	r, newState := state.nextTailRune()
	assert.Equal(t, 'o', r)
	assert.Equal(t, "hell", newState.remaining())
}

func TestStateNextHeadRuneIf(t *testing.T) {
	state := WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	r, newState := state.nextHeadRuneIf(func(r rune) bool {
		return r == 'h'
	})
	assert.Equal(t, 'h', r)
	assert.Equal(t, "ello", newState.remaining())
}

func TestStateNextTailRuneIf(t *testing.T) {
	state := WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	r, newState := state.nextTailRuneIf(func(r rune) bool {
		return r == 'o'
	})
	assert.Equal(t, 'o', r)
	assert.Equal(t, "hell", newState.remaining())
}

func TestStateConsumeTail(t *testing.T) {
	state := WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	state = state.tailConsume(2)
	assert.Equal(t, "hel", state.remaining())
}

func TestStateConsumeHead(t *testing.T) {
	state := WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	state = state.consume(2)
	assert.Equal(t, "llo", state.remaining())
}

func TestStateConsumeHeadTail(t *testing.T) {
	var state State
	state = WithState("hello")
	assert.Equal(t, "hello", state.remaining())
	state = state.consume(2)
	assert.Equal(t, "llo", state.remaining())
	state = state.tailConsume(1)
	assert.Equal(t, "ll", state.remaining())
}

func TestSnipe(t *testing.T) {
	state := WithState("banana")
	assert.Equal(t, "banana", state.remaining())
	found := state.Snipe(func(r rune) bool { return 'a' == r })
	assert.Equal(t, 3, len(found))
	assert.Equal(t, "anana", found[0].remaining())
	assert.Equal(t, "ana", found[1].remaining())
	assert.Equal(t, "a", found[2].remaining())
}

func TestSnipeEmpty(t *testing.T) {
	state := WithState("banana")
	assert.Equal(t, "banana", state.remaining())
	found := state.Snipe(func(r rune) bool { return 'q' == r })
	assert.Equal(t, 0, len(found))
}

func TestTokenizeMultiple(t *testing.T) {
	state := WithState("ban ana")
	assert.Equal(t, "ban ana", state.remaining())
	tokens := state.Tokenize(true, func(r rune) bool {
		return IsWhitespace(r)
	})

	assert.Equal(t, 2, len(tokens))
	assert.Equal(t, "ban", tokens[0].remaining())
	assert.Equal(t, " ana", tokens[1].remaining())

	assert.Equal(t, 0, tokens[0].start)
	assert.Equal(t, 3, tokens[0].end)

	assert.Equal(t, 3, tokens[1].start)
	assert.Equal(t, 7, tokens[1].end)

}

func TestTokenizeMultipleSpaces(t *testing.T) {
	state := WithState("ban   ana")
	assert.Equal(t, "ban   ana", state.remaining())
	tokens := state.Tokenize(true, func(r rune) bool {
		return IsWhitespace(r)
	})

	assert.Equal(t, 4, len(tokens))
	assert.Equal(t, "ban", tokens[0].remaining())
	assert.Equal(t, " ", tokens[1].remaining())
	assert.Equal(t, " ", tokens[2].remaining())
	assert.Equal(t, " ana", tokens[3].remaining())

	assert.Equal(t, 0, tokens[0].start)
	assert.Equal(t, 3, tokens[0].end)

	assert.Equal(t, 3, tokens[1].start)
	assert.Equal(t, 4, tokens[1].end)

	assert.Equal(t, 4, tokens[2].start)
	assert.Equal(t, 5, tokens[2].end)

	assert.Equal(t, 5, tokens[3].start)
	assert.Equal(t, 9, tokens[3].end)

}

func TestTokenizeMultipleSpacesEnd(t *testing.T) {
	state := WithState("ban   ")
	assert.Equal(t, "ban   ", state.remaining())
	tokens := state.Tokenize(false, func(r rune) bool {
		return IsWhitespace(r)
	})

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "ban", tokens[0].remaining())

	assert.Equal(t, 0, tokens[0].start)
	assert.Equal(t, 3, tokens[0].end)

}

func TestTokenizeMultipleSpacesStart(t *testing.T) {
	state := WithState("   ban")
	assert.Equal(t, "   ban", state.remaining())
	tokens := state.Tokenize(false, func(r rune) bool {
		return IsWhitespace(r)
	})

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "ban", tokens[0].remaining())

	assert.Equal(t, 3, tokens[0].start)
	assert.Equal(t, 6, tokens[0].end)

}

func TestTokenizeMultipleSpacesStartEnd(t *testing.T) {
	state := WithState("   ")
	assert.Equal(t, "   ", state.remaining())
	tokens := state.Tokenize(true, func(r rune) bool {
		return IsWhitespace(r)
	})

	assert.Equal(t, 3, len(tokens))

	assert.Equal(t, " ", tokens[0].remaining())
	assert.Equal(t, " ", tokens[1].remaining())
	assert.Equal(t, " ", tokens[2].remaining())

	assert.Equal(t, 0, tokens[0].start)
	assert.Equal(t, 1, tokens[0].end)

	assert.Equal(t, 1, tokens[1].start)
	assert.Equal(t, 2, tokens[1].end)

	assert.Equal(t, 2, tokens[2].start)
	assert.Equal(t, 3, tokens[2].end)

	tokens = state.Tokenize(false, func(r rune) bool {
		return IsWhitespace(r)
	})

	assert.Equal(t, 0, len(tokens))

}

func TestTokenizeRune(t *testing.T) {
	state := WithState("banana")
	assert.Equal(t, "banana", state.remaining())
	tokens := state.Tokenize(true, func(r rune) bool {
		return r == 'n'
	})

	assert.Equal(t, 3, len(tokens))
	assert.Equal(t, "ba", tokens[0].remaining())
	assert.Equal(t, "na", tokens[1].remaining())
	assert.Equal(t, "na", tokens[2].remaining())

	assert.Equal(t, 0, tokens[0].start)
	assert.Equal(t, 2, tokens[0].end)

	assert.Equal(t, 2, tokens[1].start)
	assert.Equal(t, 4, tokens[1].end)

	assert.Equal(t, 4, tokens[2].start)
	assert.Equal(t, 6, tokens[2].end)

}

func TestContainsRune(t *testing.T) {
	state := WithState("banana")
	assert.Equal(t, "banana", state.remaining())
	assert.True(t, state.ContainsRune('a'))
	assert.False(t, state.ContainsRune('q'))
}

func TestContainsRunes(t *testing.T) {
	state := WithState("banana")
	assert.Equal(t, "banana", state.remaining())
	assert.True(t, state.ContainsAnyRunes([]rune{'a', 'n'}))
	assert.False(t, state.ContainsAnyRunes([]rune{'q', 'r'}))
}

func TestExtend(t *testing.T) {
	state := WithState("banana")
	assert.Equal(t, "banana", state.remaining())
	tail := state.tailConsume(2)
	head := state.consume(2)
	assert.Equal(t, "nana", head.remaining())
	assert.Equal(t, "bana", tail.remaining())
	assert.Equal(t, "na", head.Extend(tail).remaining())
	assert.Equal(t, "banana", tail.Extend(state).remaining())
	assert.Equal(t, "nana", head.Extend(state).remaining())
}

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

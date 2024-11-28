package parser

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	var state State
	state = State{data: "hello", offset: 0}
	assert.Equal(t, "hello", state.remaining())
	state = state.consume(2)
	assert.Equal(t, "llo", state.remaining())
	r, newState := state.nextRune()
	assert.Equal(t, 'l', r)
	assert.Equal(t, "lo", newState.remaining())

	r, newState = newState.nextRune()
	assert.Equal(t, 'l', r)
	assert.Equal(t, "o", newState.remaining())

	r, newState = newState.nextRune()
	assert.Equal(t, 'o', r)
	assert.Equal(t, "", newState.remaining())

	r, newState = newState.nextRune()
	assert.Equal(t, utf8.RuneError, r)
	assert.Equal(t, "", newState.remaining())

}

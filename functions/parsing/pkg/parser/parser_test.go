package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFail(t *testing.T) {
	state := State{data: "hello", start: 0}
	result, newState, err := Fail[string](state)
	assert.Equal(t, "no match", err.Error())
	assert.Equal(t, state, newState)
	assert.Equal(t, "", result)
}

func TestSucceed(t *testing.T) {
	state := State{data: "hello", start: 0}
	result, newState, err := Succeed("result")(state)
	assert.Nil(t, err)
	assert.Equal(t, "result", result)
	assert.Equal(t, state, newState)
}

func TestMap(t *testing.T) {
	state := State{data: "hello", start: 0}
	parser := Map(Succeed(42), func(i int) string {
		return "result"
	})
	result, newState, err := parser(state)
	assert.Nil(t, err)
	assert.Equal(t, "result", result)
	assert.Equal(t, state, newState)
}

func TestAndThen(t *testing.T) {
	state := WithState("hello")
	parser := AndThen(Succeed(42), func(i int) Parser[string] {
		return Succeed("result")
	})
	result, newState, err := parser(state)
	assert.Nil(t, err)
	assert.Equal(t, "result", result)
	assert.Equal(t, state, newState)
}

func TestOneOf(t *testing.T) {
	state := WithState("hello")
	parser := OneOf(Succeed(42), Succeed(43))
	result, newState, err := parser(state)
	assert.Nil(t, err)
	assert.Equal(t, 42, result)
	assert.Equal(t, state, newState)
}

func TestGetString(t *testing.T) {
	state := WithState("hello")
	parser := GetString(Exactly("hello"))
	result, newState, err := parser(state)
	assert.Nil(t, err)
	assert.Equal(t, "hello", result)
	assert.Equal(t, State{data: "hello", start: 5, end: 5}, newState)
}

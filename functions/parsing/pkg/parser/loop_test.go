package parser

import (
	"testing"
)

func TestLoop(t *testing.T) {

	step := Step[int, string]{
		Done:  true,
		Accum: 42,
		Value: "foo",
	}

	stepper := func(accum int) Parser[Step[int, string]] {
		return func(state State) (Step[int, string], State, error) {
			return step, state, nil
		}
	}

	loop := Loop(0, stepper)

	state := State{data: "hello", offset: 0}
	result, newState, err := loop(state)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "foo" {
		t.Errorf("unexpected result: %v", result)
	}
	if newState.data != "hello" {
		t.Errorf("unexpected state: %v", newState)
	}
}

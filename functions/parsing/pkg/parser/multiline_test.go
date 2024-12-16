package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultilineParser(t *testing.T) {

	w1 := StartSkipping(Exactly("<"))
	k1 := AppendKeeping(w1, IntParser)
	w2 := AppendSkipping(k1, Exactly(">"))
	s2 := AppendKeeping(w2, Map(OneOf(
		NumberStringParser,
		MonthAsciiParser,
	), func(text string) string { return text }))

	predicate := Apply2(s2, func(p int, month string) bool {
		return true
	})

	mlp := MultilineParser('<', predicate)

	line1 := "<165>1 2021-09-15T14:00:00.000Z host "
	line2 := "<dave>"
	line3 := "<14>Dec 23 host "

	state := WithState(line1 + line2 + line3)
	results, state, err := mlp(state)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	assert.Equal(t, 2, len(results))

	if state.start != len(line1+line2) {
		t.Errorf("Expected state.start to be %v, got %v", len(line1+line2), state.start)
	}

	for i, line := range results {
		fmt.Printf("Line %v: %v\n", i, line)
	}

}

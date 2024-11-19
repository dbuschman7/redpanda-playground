package main

import (
	"testing"

	"dave.internal/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestTrueParser(t *testing.T) {
	r, err := parser.Parse(TrueParser, "true")
	if err != nil {
		t.Errorf("TrueParser failed: %v", err)
	}
	assert.True(t, r)

	_, err = parser.Parse(TrueParser, "false")
	if err == nil {
		t.Errorf("TrueParser failed: %v", err)
	}

	_, err = parser.Parse(TrueParser, "tru")
	if err == nil {
		t.Errorf("TrueParser failed: %v", err)
	}

}

func TestFalseParser(t *testing.T) {
	r, err := parser.Parse(FalseParser, "false")
	if err != nil {
		t.Errorf("FalseParser failed: %v", err)
	}
	assert.False(t, r)

	_, err = parser.Parse(FalseParser, "true")
	if err == nil {
		t.Errorf("FalseParser failed: %v", err)
	}

	_, err = parser.Parse(FalseParser, "fals")
	if err == nil {
		t.Errorf("FalseParser failed: %v", err)
	}
}

func TestBoolParser(t *testing.T) {
	r, err := parser.Parse(BoolParser, "true")
	if err != nil {
		t.Errorf("BoolParser failed: %v", err)
	}
	assert.True(t, r)

	r, err = parser.Parse(BoolParser, "false")
	if err != nil {
		t.Errorf("BoolParser failed: %v", err)
	}
	assert.False(t, r)

	_, err = parser.Parse(BoolParser, "tru")
	if err == nil {
		t.Errorf("BoolParser failed: %v", err)
	}
}

func TestIntParser(t *testing.T) {
	r, err := parser.Parse(IntParser, "123")
	if err != nil {
		t.Errorf("IntParser failed: %v", err)
	}
	assert.Equal(t, 123, r)

	_, err = parser.Parse(IntParser, "0123")
	if err == nil {
		t.Errorf("IntParser failed: %v", err)
	}

	_, err = parser.Parse(IntParser, "abc")
	if err == nil {
		t.Errorf("IntParser failed: %v", err)
	}
}

func TestNameParser(t *testing.T) {
	r, err := parser.Parse(NameParser, "abc")
	if err != nil {
		t.Errorf("NameParser failed: %v", err)
	}
	assert.Equal(t, "abc", r)

	_, err = parser.Parse(NameParser, "123")
	if err == nil {
		t.Errorf("NameParser failed: %v", err)
	}

	r, err = parser.Parse(NameParser, "a1b2c3")
	if err != nil {
		t.Errorf("NameParser failed: %v", err)
	}
	assert.Equal(t, "a1b2c3", r)
}

func TestWhitespaceSkipParser(t *testing.T) {
	r, err := parser.Parse(WhitespaceSkipParser, "  \n\t")
	if err != nil {
		t.Errorf("WhitespaceSkipParser failed: %v", err)
	}
	assert.IsType(t, parser.Empty{}, r)

	{
		s1 := parser.StartSkipping(WhitespaceSkipParser)
		s2 := parser.AppendKeeping(s1, NameParser)
		end := parser.Apply(s2, func(r string) string {
			return r
		})

		r, err := parser.Parse(end, "  abc")
		if err != nil {
			t.Errorf("WhitespaceSkipParser failed: %v", err)
		}
		assert.Equal(t, "abc", r)
	}
}

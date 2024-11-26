package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrueParser(t *testing.T) {
	r, err := Parse(TrueParser, "true")
	if err != nil {
		t.Errorf("TrueParser failed: %v", err)
	}
	assert.True(t, r)

	_, err = Parse(TrueParser, "false")
	if err == nil {
		t.Errorf("TrueParser failed: %v", err)
	}

	_, err = Parse(TrueParser, "tru")
	if err == nil {
		t.Errorf("TrueParser failed: %v", err)
	}

}

func TestFalseParser(t *testing.T) {
	r, err := Parse(FalseParser, "false")
	if err != nil {
		t.Errorf("FalseParser failed: %v", err)
	}
	assert.False(t, r)

	_, err = Parse(FalseParser, "true")
	if err == nil {
		t.Errorf("FalseParser failed: %v", err)
	}

	_, err = Parse(FalseParser, "fals")
	if err == nil {
		t.Errorf("FalseParser failed: %v", err)
	}
}

func TestBoolParser(t *testing.T) {
	r, err := Parse(BoolParser, "true")
	if err != nil {
		t.Errorf("BoolParser failed: %v", err)
	}
	assert.True(t, r)

	r, err = Parse(BoolParser, "false")
	if err != nil {
		t.Errorf("BoolParser failed: %v", err)
	}
	assert.False(t, r)

	_, err = Parse(BoolParser, "tru")
	if err == nil {
		t.Errorf("BoolParser failed: %v", err)
	}
}

func TestIntParser(t *testing.T) {
	r, err := Parse(IntParser, "123")
	if err != nil {
		t.Errorf("IntParser failed: %v", err)
	}
	assert.Equal(t, 123, r)

	_, err = Parse(IntParser, "0123")
	if err == nil {
		t.Errorf("IntParser failed: %v", err)
	}

	_, err = Parse(IntParser, "abc")
	if err == nil {
		t.Errorf("IntParser failed: %v", err)
	}
}

func TestWhitespaceSkipParser(t *testing.T) {
	r, err := Parse(WhitespaceSkipParser, "  \n\t")
	if err != nil {
		t.Errorf("WhitespaceSkipParser failed: %v", err)
	}
	assert.IsType(t, Empty{}, r)

	{
		s1 := StartSkipping(WhitespaceSkipParser)
		s2 := AppendKeeping(s1, IntParser)

		r, err := Parse(s2, "  123")

		if err != nil {
			t.Errorf("WhitespaceSkipParser failed: %v", err)
		}
		assert.Equal(t, 123, r.Second)
	}
}

func TestNameParser(t *testing.T) {
	r, err := Parse(NameParser, "abc")
	if err != nil {
		t.Errorf("NameParser failed: %v", err)
	}
	assert.Equal(t, "abc", r)

	_, err = Parse(NameParser, "123")
	if err == nil {
		t.Errorf("NameParser failed: %v", err)
	}

	r, err = Parse(NameParser, "a1b2c3")
	if err != nil {
		t.Errorf("NameParser failed: %v", err)
	}
	assert.Equal(t, "a1b2c3", r)
}

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuotedStringParser(t *testing.T) {
	state := WithState(`"abc"`)
	r, next, err := QuotedStringParser()(state)
	if err != nil {
		t.Errorf("QuotedStringParser failed: %v", err)
	}
	assert.Equal(t, "abc", r)
	assert.Equal(t, "", next.remaining())
	assert.Equal(t, 5, next.start)

	state = WithState(`'abc'`)
	r, next, err = QuotedStringParser()(state)
	if err != nil {
		t.Errorf("QuotedStringParser failed: %v", err)
	}
	assert.Equal(t, "abc", r)
	assert.Equal(t, "", next.remaining())
	assert.Equal(t, 5, next.start)

	state = WithState(`abc`)
	r, next, err = QuotedStringParser()(state)
	if err == nil {
		t.Errorf("QuotedStringParser failed: %v", err)
	}
	assert.Equal(t, 0, next.start)
	assert.Equal(t, "", r)

	state = WithState(`"abc`)
	r, next, err = QuotedStringParser()(state)
	if err == nil {
		t.Errorf("QuotedStringParser failed: %v", err)
	}
	assert.Equal(t, 0, next.start)
	assert.Equal(t, `"abc`, next.remaining())
	assert.Equal(t, "", r)

	_, next, err = QuotedStringParser()(WithState(`abc"`))
	if err == nil {
		t.Errorf("QuotedStringParser failed: %v", err)
	}
	assert.Equal(t, 0, next.start)
	assert.Equal(t, `abc"`, next.remaining())
}

func TestTailQuotedString(t *testing.T) {

	state := WithState("foo=\"hello world\"")

	s, next, err := TailQuotedStringParser()(state)

	if err != nil {
		t.Error(err)
	}

	if s != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", s)
	}

	if next.remaining() != "foo=" {
		t.Errorf("Expected empty string, got %s", next.remaining())
	}

}

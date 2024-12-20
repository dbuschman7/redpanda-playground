package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrueParser(t *testing.T) {
	r, err := Parse(TrueParser, WithState("true"))
	if err != nil {
		t.Errorf("TrueParser failed: %v", err)
	}
	assert.True(t, r)

	_, err = Parse(TrueParser, WithState("false"))
	if err == nil {
		t.Errorf("TrueParser failed: %v", err)
	}

	_, err = Parse(TrueParser, WithState("tru"))
	if err == nil {
		t.Errorf("TrueParser failed: %v", err)
	}

}

func TestFalseParser(t *testing.T) {
	r, err := Parse(FalseParser, WithState("false"))
	if err != nil {
		t.Errorf("FalseParser failed: %v", err)
	}
	assert.False(t, r)

	_, err = Parse(FalseParser, WithState("true"))
	if err == nil {
		t.Errorf("FalseParser failed: %v", err)
	}

	_, err = Parse(FalseParser, WithState("fals"))
	if err == nil {
		t.Errorf("FalseParser failed: %v", err)
	}
}

func TestBoolParser(t *testing.T) {
	r, err := Parse(BoolParser, WithState("true"))
	if err != nil {
		t.Errorf("BoolParser failed: %v", err)
	}
	assert.True(t, r)

	r, err = Parse(BoolParser, WithState("false"))
	if err != nil {
		t.Errorf("BoolParser failed: %v", err)
	}
	assert.False(t, r)

	_, err = Parse(BoolParser, WithState("tru"))
	if err == nil {
		t.Errorf("BoolParser failed: %v", err)
	}
}

func TestIntParser(t *testing.T) {
	r, err := Parse(IntParser, WithState("123"))
	if err != nil {
		t.Errorf("IntParser failed: %v", err)
	}
	assert.Equal(t, 123, r)

	_, err = Parse(IntParser, WithState("0123"))
	if err == nil {
		t.Errorf("IntParser failed: %v", err)
	}

	_, err = Parse(IntParser, WithState("abc"))
	if err == nil {
		t.Errorf("IntParser failed: %v", err)
	}
}

func TestWhitespaceSkipParser(t *testing.T) {
	r, err := Parse(WhitespaceSkipParser, WithState("  \n\t"))
	if err != nil {
		t.Errorf("WhitespaceSkipParser failed: %v", err)
	}
	assert.IsType(t, Empty{}, r)

	{
		s1 := StartSkipping(WhitespaceSkipParser)
		s2 := AppendKeeping(s1, IntParser)

		r, err := Parse(s2, WithState("  123"))

		if err != nil {
			t.Errorf("WhitespaceSkipParser failed: %v", err)
		}
		assert.Equal(t, 123, r.Second)
	}
}

func TestNameParser(t *testing.T) {
	r, err := Parse(EntityNameParser, WithState("abc"))
	if err != nil {
		t.Errorf("NameParser failed: %v", err)
	}
	assert.Equal(t, "abc", r)

	_, err = Parse(EntityNameParser, WithState("123"))
	if err == nil {
		t.Errorf("NameParser failed: %v", err)
	}

	r, err = Parse(EntityNameParser, WithState("a1b2c3"))
	if err != nil {
		t.Errorf("NameParser failed: %v", err)
	}
	assert.Equal(t, "a1b2c3", r)
}

func TestConsumeIf(t *testing.T) {
	state := WithState("abc")
	r, next, err := ConsumeIf(IsAsciiLetter)(state)
	if err != nil {
		t.Errorf("ConsumeIf failed: %v", err)
	}
	assert.IsType(t, Empty{}, r)
	assert.Equal(t, "bc", next.remaining())
	assert.Equal(t, 1, next.start)

	//
	state = WithState("123")
	_, next, err = ConsumeIf(IsAsciiLetter)(state)
	if err == nil {
		t.Errorf("ConsumeIf failed: %v", err)
	}
	assert.IsType(t, Empty{}, r)
	assert.Equal(t, "123", next.remaining())
	assert.Equal(t, 0, next.start)
}

func TestIpTupleParser(t *testing.T) {

	r, err := Parse[int](IpTupleParser(), WithState("234"))
	if err != nil {
		t.Errorf("IpTupleParser failed: %v", err)
	}
	assert.Equal(t, 234, r)

	_, err = Parse[int](IpTupleParser(), WithState("400"))
	if err == nil {
		t.Errorf("IpTupleParser failed: %v", err)
	}
}

func TestIpAddressParser(t *testing.T) {
	r, err := Parse(IpAddressParser(), WithState("1.12.123.234"))
	if err != nil {
		t.Errorf("IpAddressParser failed: %v", err)
	}
	assert.Equal(t, "1.12.123.234", r)

	_, err = Parse(IpAddressParser(), WithState("1,2,3,400"))
	if err == nil {
		t.Errorf("IpAddressParser failed: %v", err)
	}
}

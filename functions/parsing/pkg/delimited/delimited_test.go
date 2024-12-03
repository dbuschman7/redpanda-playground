package delimited

import (
	"testing"

	. "dave.internal/pkg/parser"

	"github.com/stretchr/testify/assert"
)

var intBoolValue = OneOf(
	Map(BoolParser,
		func(v bool) BindingValue {
			return BindingBool(v)
		}),
	Map(IntParser,
		func(i int) BindingValue {
			return BindingInt(i)
		}),
)

func TestDelimitedParserNameValueUseCase(t *testing.T) {

	k1 := StartKeeping(NameParser)
	s1 := AppendSkipping(k1, WhitespaceSkipParser)
	s2 := AppendSkipping(s1, Exactly("="))
	s3 := AppendSkipping(s2, WhitespaceSkipParser)
	k2 := AppendKeeping(s3, intBoolValue)
	bindingParser := Apply2(k2, func(name string, value BindingValue) Binding {
		return Binding{Name: name, Value: value}
	})

	p := DelimitedParser(bindingParser, rune(','))
	r, err := Parse(p, WithState("name1 = true, name2 = 123, name3 = false"))
	if err != nil {
		t.Errorf("DelimitedParser failed: %v", err)
	}

	assert.Equal(t, []Binding{
		{Name: "name3", Value: BindingBool(false)},
		{Name: "name2", Value: BindingInt(123)},
		{Name: "name1", Value: BindingBool(true)},
	}, r)
}

func TestCommaSeparatedValuesParser(t *testing.T) {

	s1 := StartSkipping(WhitespaceSkipParser)
	k1 := AppendKeeping(s1, intBoolValue)
	bindingParser := Apply(k1, func(value BindingValue) Binding {
		return Binding{Value: value}
	})

	p := DelimitedParser(bindingParser, rune(','))
	r, err := Parse(p, WithState("1, 2, true, 4, false"))
	if err != nil {
		t.Errorf("DelimitedParser failed: %v", err)
	}

	assert.Equal(t, []Binding{
		{Value: BindingBool(false)},
		{Value: BindingInt(4)},
		{Value: BindingBool(true)},
		{Value: BindingInt(2)},
		{Value: BindingInt(1)},
	}, r)
}

func TestPipeSeparatedValuesParser(t *testing.T) {

	s1 := StartSkipping(WhitespaceSkipParser)
	k1 := AppendKeeping(s1, intBoolValue)
	bindingParser := Apply(k1, func(value BindingValue) Binding {
		return Binding{Value: value}
	})

	p := DelimitedParser(bindingParser, rune('|'))
	r, err := Parse(p, WithState("1| 2| true| 4| false  "))
	if err != nil {
		t.Errorf("DelimitedParser failed: %v", err)
	}

	assert.Equal(t, []Binding{
		{Value: BindingBool(false)},
		{Value: BindingInt(4)},
		{Value: BindingBool(true)},
		{Value: BindingInt(2)},
		{Value: BindingInt(1)},
	}, r)
}

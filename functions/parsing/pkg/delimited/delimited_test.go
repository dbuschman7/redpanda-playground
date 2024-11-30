package delimited

import (
	"fmt"
	"testing"

	. "dave.internal/pkg/parser"

	"github.com/stretchr/testify/assert"
)

func TestDelimitedParserNameValueUseCase(t *testing.T) {

	intBoolValue := OneOf(
		Map(BoolParser,
			func(v bool) BindingValue {
				return BindingBool(v)
			}),
		Map(IntParser,
			func(i int) BindingValue {
				return BindingInt(i)
			}),
	)

	s0 := StartSkipping(WhitespaceSkipParser)
	k1 := AppendKeeping(s0, NameParser)
	s1 := AppendSkipping(k1, WhitespaceSkipParser)
	s2 := AppendSkipping(s1, Exactly("="))
	s3 := AppendSkipping(s2, WhitespaceSkipParser)
	k2 := AppendKeeping(s3, intBoolValue)
	bindingParser := Apply2(k2, func(name string, value BindingValue) Binding {
		return Binding{Name: name, Value: value}
	})

	p := DelimitedParser(bindingParser, ",")
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

func TestNonDelimiterParser(t *testing.T) {
	p := ConsumeWhileNotDelimiterParser(',', false)
	r, err := ParseSome(p, WithState("true"))
	if err != nil {
		t.Errorf("NotDelimiterParser failed: %v", err)
	}

	fmt.Printf("r: %v\n", r.Value)
	assert.Equal(t, BindingString("true"), r.Value)
}

func TestNonDelimiterParserUnquote(t *testing.T) {
	p := ConsumeWhileNotDelimiterParser(',', true)
	r, err := ParseSome(p, WithState("\"true\", true"))
	if err != nil {
		t.Errorf("NotDelimiterParser failed: %v", err)
	}

	assert.Equal(t, r, Binding{Value: BindingString("true")})

	p = ConsumeWhileNotDelimiterParser(',', false)
	r, err = ParseSome(p, WithState("\"true\""))
	if err != nil {
		t.Errorf("NotDelimiterParser failed: %v", err)
	}
	assert.Equal(t, r, Binding{Value: BindingString("\"true\"")})
}

func TestDelimitedParserCommaDelimited(t *testing.T) {

	p := ConsumeWhileNotDelimiterParser(',', true)
	r, err := Parse(DelimitedParser(p, ","), WithState("true, \"123\", false"))
	if err != nil {
		t.Errorf("DelimitedParser failed: %v", err)
	}

	assert.Equal(t, []Binding{
		{Value: BindingString("true")},
		{Value: BindingString("123")},
		{Value: BindingString("false")},
	}, r)
}

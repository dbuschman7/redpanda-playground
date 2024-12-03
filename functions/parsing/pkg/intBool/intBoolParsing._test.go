package intBool

import (
	"testing"

	. "dave.internal/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestValueParser(t *testing.T) {
	p := IntBoolMappingParser()
	r, err := Parse(p.valueParser, WithState("true"))
	if err != nil {
		t.Errorf("ValueParser failed: %v", err)
	}
	assert.Equal(t, r, BindingBool(true))

	r, err = Parse(p.valueParser, WithState("false"))
	if err != nil {
		t.Errorf("ValueParser failed: %v", err)
	}
	assert.Equal(t, r, BindingBool(false))

	r, err = Parse(p.valueParser, WithState("123"))
	if err != nil {
		t.Errorf("ValueParser failed: %v", err)
	}
	assert.Equal(t, r, BindingInt(123))
}

func TestBindingParser(t *testing.T) {
	p := IntBoolMappingParser()
	r, err := Parse(p.bindingParser, WithState("name = true"))
	if err != nil {
		t.Errorf("BindingParser failed: %v", err)
	}
	assert.Equal(t, r, Binding{Name: "name", Value: BindingBool(true)})

	r, err = Parse(p.bindingParser, WithState("name = 123"))
	if err != nil {
		t.Errorf("BindingParser failed: %v", err)
	}
	assert.Equal(t, r, Binding{Name: "name", Value: BindingInt(123)})
}

func TestBindingsParser(t *testing.T) {
	p := IntBoolMappingParser()
	r, err := Parse(p.bindingsParser, WithState("name1 = true, name2 = 123"))
	if err != nil {
		t.Errorf("BindingsParser failed: %v", err)
	}
	assert.Equal(t, []Binding{
		{Name: "name2", Value: BindingInt(123)},
		{Name: "name1", Value: BindingBool(true)},
	}, r)
}

func TestConfigurationParser(t *testing.T) {
	p := IntBoolMappingParser()
	r, err := Parse(p.ConfigurationParser, WithState("[name1 = true, name2 = 123 ]"))
	if err != nil {
		t.Errorf("ConfigurationParser failed: %v", err)
	}
	assert.Equal(t, []Binding{
		{Name: "name2", Value: BindingInt(123)},
		{Name: "name1", Value: BindingBool(true)},
	}, r)
}

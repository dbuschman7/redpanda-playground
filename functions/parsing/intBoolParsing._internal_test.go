package main

import (
	"testing"

	"dave.internal/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestValueParser(t *testing.T) {
	p := IntBoolParser()
	r, err := parser.Parse(p.valueParser, "true")
	if err != nil {
		t.Errorf("ValueParser failed: %v", err)
	}
	assert.Equal(t, r, parser.BindingBool(true))

	r, err = parser.Parse(p.valueParser, "false")
	if err != nil {
		t.Errorf("ValueParser failed: %v", err)
	}
	assert.Equal(t, r, parser.BindingBool(false))

	r, err = parser.Parse(p.valueParser, "123")
	if err != nil {
		t.Errorf("ValueParser failed: %v", err)
	}
	assert.Equal(t, r, parser.BindingInt(123))
}

func TestBindingParser(t *testing.T) {
	p := IntBoolParser()
	r, err := parser.Parse(p.bindingParser, "name = true")
	if err != nil {
		t.Errorf("BindingParser failed: %v", err)
	}
	assert.Equal(t, r, parser.Binding{Name: "name", Value: parser.BindingBool(true)})

	r, err = parser.Parse(p.bindingParser, "name = 123")
	if err != nil {
		t.Errorf("BindingParser failed: %v", err)
	}
	assert.Equal(t, r, parser.Binding{Name: "name", Value: parser.BindingInt(123)})
}

func TestBindingsParser(t *testing.T) {
	p := IntBoolParser()
	r, err := parser.Parse(p.bindingsParser, "name1 = true, name2 = 123")
	if err != nil {
		t.Errorf("BindingsParser failed: %v", err)
	}
	assert.Equal(t, r, []parser.Binding{
		{Name: "name2", Value: parser.BindingInt(123)},
		{Name: "name1", Value: parser.BindingBool(true)},
	})
}

func TestConfigurationParser(t *testing.T) {
	p := IntBoolParser()
	r, err := parser.Parse(p.ConfigurationParser, "[name1 = true, name2 = 123 ]")
	if err != nil {
		t.Errorf("ConfigurationParser failed: %v", err)
	}
	assert.Equal(t, r, []parser.Binding{
		{Name: "name2", Value: parser.BindingInt(123)},
		{Name: "name1", Value: parser.BindingBool(true)},
	})
}

package parser

import (
	"bytes"
	"fmt"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteBinding(t *testing.T) {
	var buffer bytes.Buffer
	b := Binding{
		Name:  "foo",
		Value: BindingInt(42),
	}
	WriteBindingValueAsJson(&buffer, b)
	assert.Equal(t, " \"foo\": 42", buffer.String())
}

func TestWriteBindingBool(t *testing.T) {
	var buffer bytes.Buffer
	b := Binding{
		Name:  "bar",
		Value: BindingBool(true),
	}
	WriteBindingValueAsJson(&buffer, b)
	assert.Equal(t, " \"bar\": true", buffer.String())
}

func TestWriteBindingString(t *testing.T) {
	var buffer bytes.Buffer
	b := Binding{
		Name:  "baz",
		Value: BindingString("false"),
	}
	WriteBindingValueAsJson(&buffer, b)
	assert.Equal(t, " \"baz\": \"false\"", buffer.String())
}

func TestWriteBindingBinding(t *testing.T) {
	var buffer bytes.Buffer

	b := Binding{
		Name: "foo",
		Value: BindingBinding{
			{
				Name:  "bar",
				Value: BindingInt(42),
			},
			{
				Name:  "baz",
				Value: BindingBool(true),
			},
		},
	}
	WriteBindingValueAsJson(&buffer, b)
	fmt.Printf("buffer: %v\n", buffer.String())
	assert.Equal(t, " \"foo\": { \"bar\": 42, \"baz\": true }", buffer.String())
}

func TestWriteBindings(t *testing.T) {

	bindings := []Binding{
		{
			Name:  "foo",
			Value: BindingInt(42),
		},
		{
			Name:  "bar",
			Value: BindingBool(true),
		},
		{
			Name: "baz",
			Value: BindingBinding{
				{
					Name:  "qux",
					Value: BindingInt(42),
				},
				{
					Name:  "quux",
					Value: BindingBool(true),
				},
			},
		},
	}

	var buffer bytes.Buffer
	WriteBindingsAsJson(&buffer, "some raw message here", bindings, nil)
	fmt.Printf("buffer: %v\n", buffer.String())

	assert.Equal(t,
		"{ \"parsed\": true, \"bindings\": { \"foo\": 42, \"bar\": true, \"baz\": { \"qux\": 42, \"quux\": true } }, \"raw\": \"some raw message here\" }",
		buffer.String())
}

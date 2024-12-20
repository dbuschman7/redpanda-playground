package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// The result of parsing is a slice of Bindings
type BindingList []Binding

type BindingTree struct {
	binding Binding
	next    *BindingTree
}

// A Binding corresponds to “name = value”
type Binding struct {
	Name  string
	Value BindingValue
}

// BindingValue is a marker interface for the values in a Binding.
type BindingValue interface {
	IsBindingValue()
}

// BindingInt is a wrapper on int to implement the BindingValue interface.
type BindingInt int

// The marker method to be a BindingValue
func (BindingInt) IsBindingValue() {}

// BindingBool is a wrapper on bool to implement the BindingValue interface.
type BindingBool bool

// The marker method to be a BindingValue
func (BindingBool) IsBindingValue() {}

// BindingString is a wrapper on bool to implement the BindingValue interface.
type BindingString string

// The marker method to be a BindingValue
func (BindingString) IsBindingValue() {}

type BindingBinding []Binding

func (BindingBinding) IsBindingValue() {}

func WriteBindingValueAsJson(buffer *bytes.Buffer, b Binding) error {

	var retErr error

	buffer.WriteString(" \"")
	buffer.WriteString(b.Name)
	buffer.WriteString("\": ")

	switch v := b.Value.(type) {
	case BindingInt:
		buffer.WriteString(strconv.Itoa(int(v)))
	case BindingBool:
		buffer.WriteString(strconv.FormatBool(bool(v)))
	case BindingString:
		enc, err := json.Marshal(string(v))
		if err != nil {
			buffer.WriteString("\"error\"")
			retErr = fmt.Errorf("error encoding string %v", string(v))
		} else {
			buffer.Write(enc)
		}
	case BindingBinding:
		buffer.WriteString("{")
		first := true
		for _, v := range v {
			if !first {
				buffer.WriteString(",")
			}
			first = false
			WriteBindingValueAsJson(buffer, v)
		}
		buffer.WriteString(" }")
	}
	return retErr
}

func WriteBindingsAsJson(buffer *bytes.Buffer, raw string, bindings []Binding, err error) {

	var document []Binding

	document = append(document, Binding{Name: "parsed", Value: BindingBool(err == nil)})
	if err != nil {
		document = append(document, Binding{Name: "error", Value: BindingString(err.Error())})
	} else {
		document = append(document, Binding{Name: "bindings", Value: BindingBinding(bindings)})
	}
	if raw != "" {
		document = append(document, Binding{Name: "raw", Value: BindingString(raw)})
	}

	buffer.WriteString("{")
	first := true
	for _, b := range document {
		if !first {
			buffer.WriteString(",")
		}
		first = false
		WriteBindingValueAsJson(buffer, b) // TODO handle error from Json encoding
	}
	buffer.WriteString(" }")

}

var bindingNameParser = OneOf(StringParser, QuotedStringParser())

// non-recursively defined parser for a BindingValue
var bindingValueParserValue = OneOf(
	Map(BoolParser,
		func(v bool) BindingValue {
			return BindingBool(v)
		}),
	Map(IntParser,
		func(i int) BindingValue {
			return BindingInt(i)
		}),
	Map(QuotedStringParser(),
		func(s string) BindingValue {
			return BindingString(s)
		}),
	Map(StringParser,
		func(s string) BindingValue {
			return BindingString(s)
		}),
)

func BindingParser() Parser[Binding] {
	s := StartKeeping(bindingNameParser)
	s1 := AppendSkipping(s, WhitespaceSkipParser)
	s2 := AppendSkipping(s1, Exactly("="))
	s3 := AppendSkipping(s2, WhitespaceSkipParser)
	s4 := AppendKeeping(s3, bindingValueParserValue)

	return Apply2(s4,
		func(name string, value BindingValue) Binding {
			return Binding{Name: name, Value: value}
		})
}

func BindingsParser(delimiter rune) Parser[BindingList] {
	return DelimitedParser(BindingParser(), delimiter)
}

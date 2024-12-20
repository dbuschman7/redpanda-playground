package intBool

import (
	"dave.internal/pkg/parser"
)

// This package provides a Parser that parses a simple configuration file format
// made up as an example.  Here is a grammar for the configuration format:
//
//	configuration:  '[' whitespace bindings whitespace ']'
//
//	bindings: binding | binding whitespace ',' whitespace bindings
//
//	binding:  name whitespace '=' whitespace value
//
//	name: [a-zA-Z][0-9a-zA-Z]*
//
//	value:  int | bool
//
//	int: [0-9] | [1-9][0-9]+
//
//	bool: "true" | "false"
//
//	whitespace: [ \t\n]*
//
// IntBoolMappingParser returns a ConfigParsers structure containing the
// ConfigurationParser that you actually want to pass to Parse
// if you want to parse the Configuration format.
// NewConfigParser returns this struct.  The sole exported field is the Parser for the
// entire configuration format.  The unexported fields contain subcomponent parsers
// for mutual reference and (ideally) internal testing.
type IntBoolMappingParsers struct {
	trueParser          parser.Parser[bool]
	falseParser         parser.Parser[bool]
	boolParser          parser.Parser[bool]
	intParser           parser.Parser[int]
	valueParser         parser.Parser[parser.BindingValue]
	nameParser          parser.Parser[string]
	bindingParser       parser.Parser[parser.Binding]
	whitespaceParser    parser.Parser[parser.Empty]
	bindingsParser      parser.Parser[parser.BindingList]
	ConfigurationParser parser.Parser[parser.BindingList]
}

func IntBoolMappingParser() IntBoolMappingParsers {
	var p IntBoolMappingParsers

	p.trueParser = parser.TrueParser
	p.falseParser = parser.FalseParser
	p.boolParser = parser.BoolParser
	p.intParser = parser.IntParser

	p.valueParser = parser.OneOf(
		parser.Map(p.boolParser,
			func(v bool) parser.BindingValue {
				return parser.BindingBool(v)
			}),
		parser.Map(p.intParser,
			func(i int) parser.BindingValue {
				return parser.BindingInt(i)
			}),
	)

	p.nameParser = parser.EntityNameParser
	p.whitespaceParser = parser.WhitespaceSkipParser

	{
		s := parser.StartKeeping(p.nameParser)
		s1 := parser.AppendSkipping(s, p.whitespaceParser)
		s2 := parser.AppendSkipping(s1, parser.Exactly("="))
		s3 := parser.AppendSkipping(s2, p.whitespaceParser)
		s4 := parser.AppendKeeping(s3, p.valueParser)
		p.bindingParser = parser.Apply2(s4,
			func(name string, value parser.BindingValue) parser.Binding {
				return parser.Binding{Name: name, Value: value}
			})
	}
	{
		p.bindingsParser = parser.DelimitedParser(p.bindingParser, ',')
	}
	{
		s := parser.StartSkipping(parser.Exactly("["))
		s1 := parser.AppendSkipping(s, p.whitespaceParser)
		s2 := parser.AppendKeeping(s1, p.bindingsParser)
		s3 := parser.AppendSkipping(s2, p.whitespaceParser)
		s4 := parser.AppendSkipping(s3, parser.Exactly("]"))
		p.ConfigurationParser = parser.Apply(s4, func(b parser.BindingList) parser.BindingList { return b })
	}
	return p
}

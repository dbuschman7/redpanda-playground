package main

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
// IntBoolParser returns a ConfigParsers structure containing the
// ConfigurationParser that you actually want to pass to Parse
// if you want to parse the Configuration format.
// NewConfigParser returns this struct.  The sole exported field is the Parser for the
// entire configuration format.  The unexported fields contain subcomponent parsers
// for mutual reference and (ideally) internal testing.
type ConfigParsers struct {
	trueParser          parser.Parser[bool]
	falseParser         parser.Parser[bool]
	boolParser          parser.Parser[bool]
	intParser           parser.Parser[int]
	valueParser         parser.Parser[parser.BindingValue]
	nameParser          parser.Parser[string]
	bindingParser       parser.Parser[parser.Binding]
	whitespaceParser    parser.Parser[parser.Empty]
	bindingsParser      parser.Parser[[]parser.Binding]
	ConfigurationParser parser.Parser[[]parser.Binding]
}

func IntBoolParser() ConfigParsers {
	var p ConfigParsers

	p.trueParser = TrueParser
	p.falseParser = FalseParser
	p.boolParser = BoolParser
	p.intParser = IntParser

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

	p.nameParser = NameParser
	p.whitespaceParser = WhitespaceSkipParser

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
		type BindingList struct {
			binding parser.Binding
			next    *BindingList
		}

		p.bindingsParser = parser.Loop(nil,
			func(bindings *BindingList) parser.Parser[parser.Step[*BindingList, []parser.Binding]] {
				if bindings == nil {
					return parser.Map(p.bindingParser,
						func(binding parser.Binding) parser.Step[*BindingList, []parser.Binding] {
							return parser.Step[*BindingList, []parser.Binding]{Accum: &BindingList{binding: binding}, Done: false}
						},
					)
				}
				s := parser.StartSkipping(p.whitespaceParser)
				s1 := parser.AppendSkipping(s, parser.Exactly(","))
				s2 := parser.AppendSkipping(s1, p.whitespaceParser)
				s3 := parser.AppendKeeping(s2, p.bindingParser)
				extend := parser.Apply(s3, func(b parser.Binding) parser.Step[*BindingList, []parser.Binding] {
					return parser.Step[*BindingList, []parser.Binding]{
						Accum: &BindingList{binding: b, next: bindings},
						Done:  false,
					}
				})

				var bindingSlice []parser.Binding
				b := bindings
				for {
					if b == nil {
						break
					}
					bindingSlice = append(bindingSlice, b.binding)
					b = b.next
				}
				return parser.OneOf(
					extend,
					parser.Succeed(parser.Step[*BindingList, []parser.Binding]{Value: bindingSlice, Done: true}),
				)

			},
		)
	}
	{
		s := parser.StartSkipping(parser.Exactly("["))
		s1 := parser.AppendSkipping(s, p.whitespaceParser)
		s2 := parser.AppendKeeping(s1, p.bindingsParser)
		s3 := parser.AppendSkipping(s2, p.whitespaceParser)
		s4 := parser.AppendSkipping(s3, parser.Exactly("]"))
		p.ConfigurationParser = parser.Apply(s4, func(b []parser.Binding) []parser.Binding { return b })
	}
	return p
}

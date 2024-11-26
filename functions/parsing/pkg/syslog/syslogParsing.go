package syslog

import (
	"dave.internal/pkg/parser"
)

// RFC 3164 -
//   <PRIVAL>TIMESTAMP HOSTNAME TAG: MESSAGE
//   <13>Oct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message.

type SyslogParsers struct {
	stringParser   parser.Parser[string]
	nameParser     parser.Parser[string]
	valueParser    parser.Parser[parser.BindingValue]
	priorityParser parser.Parser[int]
	// bindingParser parser.Parser[parser.Binding]
}

// no semantical meaning, just a string
// func datePart8601() parser.Parser[string] { // Date Parser 8601 - 2021-10-22T12:34:56
// 	s := parser.StartKeeping(IntLength4Parser)
// 	s1 := parser.AppendSkipping(s, parser.Exactly("-"))
// 	s2 := parser.AppendKeeping(s1, IntLength2Parser)
// 	s3 := parser.AppendSkipping(s2, parser.Exactly("-"))
// 	s4 := parser.AppendKeeping(s3, IntLength2Parser)

// 	return parser.Apply3(s4, func(year int, month int, day int) string {
// 		return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
// 	})
// }

func SyslogParser() SyslogParsers {
	var p SyslogParsers

	p.stringParser = parser.StringParser
	p.nameParser = parser.NameParser

	p.valueParser = parser.OneOf(
		parser.Map(p.stringParser,
			func(v string) parser.BindingValue {
				return parser.BindingString(v)
			}),
		parser.Map(parser.IntParser,
			func(i int) parser.BindingValue {
				return parser.BindingInt(i)
			}),
	)

	{ // Priority Parser
		s := parser.StartSkipping(parser.Exactly("<"))
		s1 := parser.AppendKeeping(s, parser.IntParser)
		s2 := parser.AppendSkipping(s1, parser.Exactly(">"))

		p.priorityParser = parser.Apply(s2, func(p int) int {
			return p
		})

	}

	return p
}

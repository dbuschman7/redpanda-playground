package common

import (
	"fmt"
	"slices"
	"strings"

	"dave.internal/pkg/parser"
)

// TODO - only pull 2 characters
var IntLength2Parser = parser.AndThen(
	parser.IntParser,
	func(v int) parser.Parser[int] {
		if v < 0 || v > 99 {
			return parser.Fail[int]
		}
		return parser.Succeed(v)
	},
)

var IntLength4Parser = parser.AndThen(
	parser.IntParser,
	func(v int) parser.Parser[int] {
		if v < 0 || v > 9999 {
			return parser.Fail[int]
		}
		return parser.Succeed(v)
	},
)

// Time Parser - 12:34:56 (no bounds checking)
func TimeColonParser() parser.Parser[string] {
	w := parser.StartSkipping(parser.WhitespaceSkipParser)
	s := parser.AppendKeeping(w, IntLength2Parser)
	s1 := parser.AppendSkipping(s, parser.Exactly(":"))
	s2 := parser.AppendKeeping(s1, IntLength2Parser)
	s3 := parser.AppendSkipping(s2, parser.Exactly(":"))
	s4 := parser.AppendKeeping(s3, IntLength2Parser)

	return parser.Apply3(s4, func(h, m, s int) string {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	})
}

func DateYYMMDDDashParser() parser.Parser[string] {
	w := parser.StartSkipping(parser.WhitespaceSkipParser)
	s := parser.AppendKeeping(w, IntLength2Parser)
	s1 := parser.AppendSkipping(s, parser.Exactly("-"))
	s2 := parser.AppendKeeping(s1, IntLength2Parser)
	s3 := parser.AppendSkipping(s2, parser.Exactly("-"))
	s4 := parser.AppendKeeping(s3, IntLength2Parser)

	return parser.Apply3(s4, func(y, m, d int) string {
		return fmt.Sprintf("%02d-%02d-%02d", y, m, d)
	})
}

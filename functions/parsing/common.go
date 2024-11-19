package main

import (
	"strconv"

	"dave.internal/pkg/parser"
)

func isAsciiLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlphaNum(r rune) bool {
	return isAsciiLetter(r) || isDecimalDigit(r)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

var TrueParser = parser.Map(
	parser.Exactly("true"),
	func(parser.Empty) bool {
		return true
	})

var FalseParser = parser.Map(
	parser.Exactly("false"),
	func(parser.Empty) bool { return false })

var BoolParser = parser.OneOf(
	TrueParser,
	FalseParser)

var IntParser = parser.AndThen(parser.GetString(parser.ConsumeSome(isDecimalDigit)),
	func(digits string) parser.Parser[int] {
		if len(digits) > 1 && digits[0] == '0' {
			return parser.Fail[int]
		}
		v, err := strconv.Atoi(digits)
		if err != nil {
			return parser.Fail[int]
		}
		return parser.Succeed(v)
	},
)

var NameParser = parser.GetString(
	parser.AndThen(
		parser.ConsumeIf(isAsciiLetter),
		func(parser.Empty) parser.Parser[parser.Empty] {
			return parser.ConsumeWhile(isAlphaNum)
		},
	))

var WhitespaceSkipParser = parser.ConsumeWhile(isWhitespace)

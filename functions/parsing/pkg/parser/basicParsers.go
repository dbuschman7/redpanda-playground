package parser

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"slices"
)

func IsAsciiLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func IsDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsAlphaNum(r rune) bool {
	return IsAsciiLetter(r) || IsDecimalDigit(r)
}

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

var TrueParser = Map(
	Exactly("true"),
	func(Empty) bool {
		return true
	})

var FalseParser = Map(
	Exactly("false"),
	func(Empty) bool { return false })

var BoolParser = OneOf(
	TrueParser,
	FalseParser)

var IntParser = AndThen(GetString(ConsumeSome(IsDecimalDigit)),
	func(digits string) Parser[int] {
		if len(digits) > 1 && digits[0] == '0' {
			return Fail[int]
		}
		v, err := strconv.Atoi(digits)
		if err != nil {
			return Fail[int]
		}
		return Succeed(v)
	},
)

var NumberStringParser = AndThen(
	GetString(ConsumeSome(IsDecimalDigit)),
	func(s string) Parser[string] {
		//	fmt.Printf("nsp: %s\n", s)
		return Succeed(s)
	},
)

var StringParser = GetString(ConsumeSome(IsAlphaNum))
var AsciiParser = GetString(ConsumeSome(IsAsciiLetter))
var WhitespaceSkipParser = ConsumeWhile(IsWhitespace)

var NameParser = GetString(
	AndThen(
		ConsumeIf(IsAsciiLetter),
		func(Empty) Parser[Empty] {
			return ConsumeWhile(IsAlphaNum)
		},
	))

var quotesList = []rune{'"', '\'', '“', '”', '‘', '’', '`'}

func QuotedStringParser() Parser[string] {
	return func(initial State) (string, State, error) {
		quote, _ := utf8.DecodeRuneInString(initial.remaining())

		if !slices.Contains(quotesList, quote) {
			return "", initial, ErrNoMatch
		}

		start := initial.offset + 1
		current := start
		found := false
		for pos, char := range initial.data[start:] {
			if char == quote {
				found = true
				return initial.data[start:current], initial.consume(pos + 1), nil
			}
			if char == '\\' {
				current = current + 1 // skip an extra character
			}
			current = current + 1
		}

		if !found {
			return "", initial, fmt.Errorf("no closing quote")
		}
		return "", initial, nil

	}
}

func EitherOrParser(first string, second string) Parser[string] {
	return OneOf(
		Map(Exactly(first), func(Empty) string { return first }),
		Map(Exactly(second), func(Empty) string { return second }),
	)
}

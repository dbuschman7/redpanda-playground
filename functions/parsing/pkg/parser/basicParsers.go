package parser

import (
	"fmt"
	"strconv"
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

func EitherOrParser(first string, second string) Parser[string] {
	return OneOf(
		Map(Exactly(first), func(Empty) string { return first }),
		Map(Exactly(second), func(Empty) string { return second }),
	)
}

func IpTupleParser() Parser[int] {
	return AndThen(
		IntParser,
		func(first int) Parser[int] {
			if first < 0 || first > 255 {
				return Fail[int]
			}
			return Succeed(first)
		})
}

func IpAddressParser() Parser[string] {
	t1 := StartKeeping(IpTupleParser())
	s2 := AppendSkipping(t1, Exactly("."))
	t2 := AppendKeeping(s2, IpTupleParser())
	s3 := AppendSkipping(t2, Exactly("."))
	t3 := AppendKeeping(s3, IpTupleParser())
	s4 := AppendSkipping(t3, Exactly("."))
	t4 := AppendKeeping(s4, IpTupleParser())

	return Apply4(t4, func(a, b, c, d int) string {
		return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
	})
}

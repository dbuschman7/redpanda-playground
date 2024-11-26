package datetime

import (
	"fmt"
	"slices"
	"strings"

	. "dave.internal/pkg/parser"
)

func isNumOrDash(r rune) bool {
	return IsDecimalDigit(r) || r == '-'
}

func isNumOrColonOrPeriod(r rune) bool {
	return IsDecimalDigit(r) || r == ':' || r == '.'
}

func isPlusMinus(r rune) bool {
	return r == '+' || r == '-'
}

func isNumOrColon(r rune) bool {
	return r == ':' || IsDecimalDigit(r)
}

func split(r rune) bool {
	return r == ':' || r == '.' || r == '-' || r == 'T' || r == ' ' || r == '+'
}

func generatePattern(text string) int {
	pattern := 0
	for _, part := range strings.FieldsFunc(text, split) {
		pattern = pattern*10 + len(part)
	}
	return pattern
}

var supportedTimeZonePatterns = []int{
	1,  // Z
	2,  // +01
	22, // +01:00
	4,  // -0100
}

func timeZone8061Support() Parser[string] {
	s := StartSkipping(WhitespaceSkipParser)

	s1 := AppendKeeping(s, OneOf[string](
		Map(Exactly("Z"), func(Empty) string { return "Z" }),
		AndThen(
			Map(GetString(ConsumeSome(isPlusMinus)), func(sign string) string { return sign }),
			func(test string) Parser[string] {
				return Map(GetString(ConsumeSome(isNumOrColon)), func(second string) string { return test + second })
			},
		),
	))

	return func(initial State) (string, State, error) {
		text, next, err := s1(initial)

		if err != nil {
			var zero string
			return zero, initial, err
		}

		pattern := generatePattern(text.Second)
		// fmt.Printf("pattern: %v - %v\n", pattern, text.Second)

		if slices.Contains(supportedTimeZonePatterns, pattern) {
			return text.Second, next, nil
		} else {
			var zero string
			return zero, initial, fmt.Errorf("unrecognized date pattern %v", pattern)
		}

	}
}

var supportedTimePatterns = []int{
	222,  // 01:01:01
	2223, // 01:01:01.000
	6,    // 010101
	63,   // 010101.000
}

func time8601Support() Parser[string] {
	s := StartSkipping(WhitespaceSkipParser)
	s1 := AppendKeeping(s, GetString(ConsumeSome(isNumOrColonOrPeriod)))
	s2 := Apply(s1, func(time string) string {
		return time
	})

	return func(initial State) (string, State, error) {
		text, next, err := s2(initial)

		if err != nil {
			var zero string
			return zero, initial, err
		}

		pattern := generatePattern(text)
		// fmt.Printf("pattern: %v - %v\n", pattern, text)
		if slices.Contains(supportedTimePatterns, pattern) {
			return text, next, nil
		} else {
			var zero string
			return zero, initial, fmt.Errorf("unrecognized date pattern %v", pattern)
		}
	}
}

var supportedDatePatterns = []int{
	422, // 2020-01-01
	8,   // 20200101
	432, // 2020-123
}

func date8601Support() Parser[string] {
	s := StartSkipping(WhitespaceSkipParser)
	s1 := AppendKeeping(s, Map(GetString(ConsumeSome(isNumOrDash)), func(text string) string { return text }))

	s2 := Apply(s1, func(date string) string {
		fmt.Printf("date: '%v'\n", date)
		return date
	})

	return func(initial State) (string, State, error) {
		text, next, err := s2(initial)

		if err != nil {
			var zero string
			return zero, initial, err
		}

		// fmt.Printf("text: '%v'\n", text)
		pattern := generatePattern(text)
		// fmt.Printf("pattern: %v\n", pattern)
		if slices.Contains(supportedDatePatterns, pattern) {
			return text, next, nil
		} else {
			var zero string
			return zero, initial, fmt.Errorf("unrecognized date pattern %v", pattern)
		}
	}
}

func ISO8601Parser() Parser[string] {
	s := StartSkipping(WhitespaceSkipParser)
	s1 := AppendKeeping(s, date8601Support())
	s2 := AppendSkipping(s1, OneOf(Exactly("T"), Exactly(" ")))
	s3 := AppendKeeping(s2, time8601Support())
	s4 := AppendKeeping(s3, timeZone8061Support())
	s5 := Apply3(s4, func(date string, time string, zone string) string {
		fmt.Printf("%s|%s|%s", date, time, zone)

		patDt := generatePattern(date)
		patTm := generatePattern(time)
		patZn := generatePattern(zone)

		if slices.Contains(supportedDatePatterns, patDt) &&
			slices.Contains(supportedTimePatterns, patTm) &&
			slices.Contains(supportedTimeZonePatterns, patZn) {
			return fmt.Sprintf("%s %s%s", date, time, zone)
		}

		return ""
	})

	return func(initial State) (string, State, error) {
		text, next, err := s5(initial)

		if err != nil {
			var zero string
			return zero, initial, err
		}

		if text == "" {
			return text, initial, fmt.Errorf("unrecognized date pattern")
		}

		return text, next, nil
	}
}

var monthListLower = []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}

var monthAsciiParser = AndThen(
	AsciiParser,
	func(month string) Parser[string] {
		if len(month) != 3 {
			return Fail[string]
		}
		found := slices.Contains(monthListLower, strings.ToLower(month)) // lookup in lowercase
		if !found {
			return Fail[string]
		}
		return Succeed(month)
	},
)

func Syslog3164Parser() Parser[string] {
	w := StartSkipping(WhitespaceSkipParser)
	s1 := AppendKeeping(w, monthAsciiParser)
	s2 := AppendSkipping(s1, WhitespaceSkipParser)
	s3 := AppendKeeping(s2, time8601Support)

	return Apply2(s3, func(month string, day int, time string) string {
		return fmt.Sprintf("%s %02d %s", month, day, time)
	})
}

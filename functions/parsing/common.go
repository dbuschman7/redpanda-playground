package main

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

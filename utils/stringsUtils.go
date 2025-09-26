package utils

import "unicode"

func IsAllUpper(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func IsValidFieldName(s string) bool {
	allowed := map[rune]bool{
		'!': true, '#': true, '$': true, '%': true,
		'&': true, '\'': true, '*': true, '+': true,
		'-': true, '.': true, '^': true, '_': true,
		'`': true, '|': true, '~': true,
	}

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || allowed[r] {
			continue
		}
		return false
	}
	return true
}

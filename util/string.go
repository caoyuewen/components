package util

import (
	"regexp"
	"unicode"
)

func HasChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func RemoveSpaces(input string) string {
	pattern := `[ \t\n\r]+`
	regex := regexp.MustCompile(pattern)
	cleaned := regex.ReplaceAllString(input, "")
	return cleaned
}

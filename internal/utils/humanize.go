package utils

import (
	"strings"
	"unicode"
)

func HumanizeName(name string) string {
	var result []rune
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) && !(unicode.IsUpper(r) && unicode.IsUpper(rune(name[i-1]))) {
			result = append(result, ' ')
		}
		result = append(result, r)
	}
	return strings.TrimSpace(string(result))
}

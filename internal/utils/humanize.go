package utils

import (
	"strings"
	"unicode"
)

func HumanizeName(name string) string {
	var result []rune
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			if (i+1 < len(name) && unicode.IsLower(rune(name[i+1]))) || (unicode.IsLower(rune(name[i-1]))) {
				result = append(result, ' ')
			}
		}

		if i == 0 || (i > 0 && result[len(result)-1] == ' ') {
			r = unicode.ToUpper(r)
		}

		result = append(result, r)
	}
	return strings.TrimSpace(string(result))
}

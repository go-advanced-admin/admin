package tests

import (
	"github.com/go-advanced-admin/admin/internal/utils"
	"testing"
)

func TestISURLSafe(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Safe", "SafeName", true},
		{"Space", "Unsafe Name", false},
		{"Safe with Special Chars", "Name_With-Special.Characters", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := utils.IsURLSafe(tt.input); result != tt.expected {
				t.Errorf("%s: IsURLSafe(%q) = %v; expected %v", tt.name, tt.input, result, tt.expected)
			}
		})
	}
}

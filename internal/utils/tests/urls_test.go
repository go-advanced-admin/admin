package tests

import (
	"github.com/go-advanced-admin/admin/internal/utils"
	"testing"
)

func TestISURLSafe(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"SafeName", true},
		{"Unsafe Name", false},
		{"Name_With-Special.Characters", true},
	}

	for _, tt := range tests {
		if result := utils.IsURLSafe(tt.input); result != tt.expected {
			t.Errorf("IsURLSafe(%q) = %v; expected %v", tt.input, result, tt.expected)
		}
	}
}

package tests

import (
	"github.com/go-advanced-admin/admin/internal/utils"
	"testing"
)

func TestHumanizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Spaces", "TestName", "Test Name"},
		{"Capitalize Start", "testName", "Test Name"},
		{"Constant", "Test", "Test"},
		{"Capitalized Word", "HTTPStatus", "HTTP Status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := utils.HumanizeName(tt.input); result != tt.expected {
				t.Errorf("HumanizeName(%s) = %s; expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

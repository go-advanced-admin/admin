package tests

import (
	"github.com/go-advanced-admin/admin/internal/utils"
	"testing"
)

func TestHumanizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"TestName", "Test Name"},
		{"testName", "Test Name"},
		{"Test", "Test"},
		{"HTTPStatus", "HTTP Status"},
	}

	for _, tt := range tests {
		if result := utils.HumanizeName(tt.input); result != tt.expected {
			t.Errorf("HumanizeName(%s) = %s; expected %s", tt.input, result, tt.expected)
		}
	}
}

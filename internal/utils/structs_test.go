package utils

import (
	"testing"
)

type TestStruct struct {
	Name  string
	Age   int
	Email string
}

func TestGetFieldValue(t *testing.T) {
	tests := []struct {
		name        string
		instance    interface{}
		fieldName   string
		expected    interface{}
		expectError bool
	}{
		{
			name:        "Valid Field Name",
			instance:    TestStruct{"John", 30, "john@example.com"},
			fieldName:   "Name",
			expected:    "John",
			expectError: false,
		},
		{
			name:        "Invalid Field Name",
			instance:    TestStruct{"John", 30, "john@example.com"},
			fieldName:   "InvalidField",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Pointer to Struct",
			instance:    &TestStruct{"Jane", 25, "jane@example.com"},
			fieldName:   "Email",
			expected:    "jane@example.com",
			expectError: false,
		},
		{
			name:        "Non-Struct Type",
			instance:    42,
			fieldName:   "Name",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Empty Struct",
			instance:    struct{}{},
			fieldName:   "Name",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := GetFieldValue(tt.instance, tt.fieldName)
			if (err != nil) != tt.expectError {
				t.Errorf("%s: unexpected error status: got %v, want %v", tt.name, err != nil, tt.expectError)
				return
			}
			if value != tt.expected {
				t.Errorf("%s: GetFieldValue(%s, %s) = %s; expected %s", tt.name, tt.instance, tt.fieldName, value, tt.expected)
			}
		})
	}
}

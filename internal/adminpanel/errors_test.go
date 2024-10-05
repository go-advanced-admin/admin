package adminpanel

import (
	"errors"
	"testing"
)

func TestGetErrorHTML(t *testing.T) {
	tests := []struct {
		name          string
		code          uint
		err           error
		expectedCode  uint
		expectedError string
	}{
		{
			name:          "Standard Error",
			code:          404,
			err:           errors.New("not found"),
			expectedCode:  404,
			expectedError: "Code: 404. Error: not found",
		},
		{
			name:          "Unauthorized Error",
			code:          401,
			err:           errors.New("unauthorized"),
			expectedCode:  401,
			expectedError: "Code: 401. Error: unauthorized",
		},
		{
			name:          "ServerError",
			code:          500,
			err:           nil,
			expectedCode:  500,
			expectedError: "Code: 500. Error: <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _ := GetErrorHTML(tt.code, tt.err)
			if code != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, code)
			}
		})
	}
}

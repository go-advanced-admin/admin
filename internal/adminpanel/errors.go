package adminpanel

import (
	"fmt"
)

// GetErrorHTML generates an HTML string representing an error message with the given code and error.
func GetErrorHTML(code uint, err error) (uint, string) {
	if err == nil {
		return code, fmt.Sprintf("Code: %v.", code)
	}
	return code, fmt.Sprintf("Code: %v. Error: %v", code, err.Error())
}

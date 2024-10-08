package adminpanel

import (
	"fmt"
)

func GetErrorHTML(code uint, err error) (uint, string) {
	if err == nil {
		return code, fmt.Sprintf("Code: %v.", code)
	}
	return code, fmt.Sprintf("Code: %v. Error: %v", code, err.Error())
}

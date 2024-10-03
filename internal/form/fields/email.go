package fields

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"net/mail"
	"strings"
)

type EmailField struct {
	BaseField
	Required bool
}

func (f *EmailField) HTML() (string, error) {
	attributesMap := make(map[string]*string)

	if f.InitialValue != nil {
		htmlType, err := f.GoTypeToHTMLType(f.InitialValue)
		if err != nil {
			return "", err
		}
		value := template.HTMLEscapeString(string(htmlType))
		attributesMap["value"] = &value
	}

	if f.Required {
		attributesMap["required"] = nil
	}

	inputType := "email"
	attributesMap["type"] = &inputType
	name := template.HTMLEscapeString(f.Name)
	attributesMap["name"] = &name

	if f.SupersedingAttributes != nil {
		for key, value := range f.SupersedingAttributes {
			attributesMap[key] = value
		}
	}

	var attributes []string
	for key, value := range attributesMap {
		if value == nil {
			attributes = append(attributes, key)
		} else {
			attributes = append(attributes, fmt.Sprintf(`%s="%s"`, key, template.HTMLEscapeString(*value)))
		}
	}

	return fmt.Sprintf(`<input %s>`, strings.Join(attributes, " ")), nil
}

func (f *EmailField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == nil {
		return "", nil
	}
	strValue, ok := value.(string)
	if !ok {
		return "", errors.New("value must be a string")
	}
	return form.HTMLType(strValue), nil
}

func (f *EmailField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	return string(value), nil
}

func (f *EmailField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()
	baseValidations = append(baseValidations, f.requiredValidation, f.emailValidation)
	return baseValidations
}

func (f *EmailField) requiredValidation(value interface{}) ([]error, error) {
	if f.Required {
		if value == nil {
			return []error{errors.New("field is required")}, nil
		}
		strValue, ok := value.(string)
		if !ok {
			return nil, errors.New("value must be a string")
		}
		if strings.TrimSpace(strValue) == "" {
			return []error{errors.New("field is required")}, nil
		}
	}
	return nil, nil
}

func (f *EmailField) emailValidation(value interface{}) ([]error, error) {
	if value == nil {
		return nil, nil
	}
	strValue, ok := value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	if _, err := mail.ParseAddress(strValue); err != nil {
		return []error{errors.New("invalid email address")}, nil
	}
	return nil, nil
}

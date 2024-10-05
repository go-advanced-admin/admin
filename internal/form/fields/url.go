package fields

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"net/url"
	"strings"
)

type URLField struct {
	BaseField
	Required bool
}

func (f *URLField) HTML() (string, error) {
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

	inputType := "url"
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

func (f *URLField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == nil {
		return "", nil
	}
	strValue, ok := value.(string)
	if !ok {
		return "", errors.New("value must be a string")
	}
	return form.HTMLType(strValue), nil
}

func (f *URLField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	return string(value), nil
}

func (f *URLField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()
	baseValidations = append(baseValidations, f.requiredValidation, f.urlValidation)
	return baseValidations
}

func (f *URLField) requiredValidation(value interface{}) ([]error, error) {
	if !f.Required {
		return nil, nil
	}
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
	return nil, nil
}

func (f *URLField) urlValidation(value interface{}) ([]error, error) {
	if value == nil {
		return nil, nil
	}
	strValue, ok := value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	_, err := url.ParseRequestURI(strValue)
	if err != nil {
		return []error{errors.New("invalid URL")}, nil
	}
	return nil, nil
}

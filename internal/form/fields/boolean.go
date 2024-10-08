package fields

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"strings"
)

type BooleanField struct {
	BaseField
	Required bool
}

func (f *BooleanField) HTML() (string, error) {
	attributesMap := make(map[string]*string)

	if f.Required {
		attributesMap["required"] = nil
	}

	inputType := "checkbox"
	attributesMap["type"] = &inputType
	name := template.HTMLEscapeString(f.Name)
	attributesMap["name"] = &name

	if f.InitialValue != nil {
		checked, ok := f.InitialValue.(bool)
		if !ok {
			return "", errors.New("initial value must be a boolean")
		}
		if checked {
			attributesMap["checked"] = nil
		}
	}

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

func (f *BooleanField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == nil {
		return "", nil
	}
	boolValue, ok := value.(bool)
	if !ok {
		return "", errors.New("value must be a boolean")
	}
	if boolValue {
		return "on", nil
	}
	return "", nil
}

func (f *BooleanField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	if string(value) == "on" || string(value) == "true" {
		return true, nil
	}
	return false, nil
}

func (f *BooleanField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()
	baseValidations = append(baseValidations, f.requiredValidation)
	return baseValidations
}

func (f *BooleanField) requiredValidation(value interface{}) ([]error, error) {
	if value == nil && f.Required {
		return []error{errors.New("field is required")}, nil
	}
	return nil, nil
}

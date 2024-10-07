package fields

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"github.com/google/uuid"
	"html/template"
	"strings"
)

type UUIDField struct {
	BaseField
	Required bool
}

func (f *UUIDField) HTML() (string, error) {
	attributesMap := make(map[string]*string)
	if f.InitialValue != nil {
		htmlType, err := f.GoTypeToHTMLType(f.InitialValue)
		if err != nil {
			return "", err
		}
		value := template.HTMLEscapeString(string(htmlType))
		attributesMap["value"] = &value
	}
	value := "36"
	attributesMap["maxLength"] = &value
	attributesMap["minLength"] = &value
	if f.Required {
		attributesMap["required"] = nil
	}
	value = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
	attributesMap["pattern"] = &value
	value = "text"
	attributesMap["type"] = &value
	value = template.HTMLEscapeString(f.Name)
	attributesMap["name"] = &value

	if f.SupersedingAttributes != nil {
		for key, value := range f.SupersedingAttributes {
			attributesMap[key] = value
		}
	}

	attributes := make([]string, 0)
	for key, value := range attributesMap {
		if value == nil {
			attributes = append(attributes, key)
		} else {
			attributes = append(attributes, fmt.Sprintf(`%s="%s"`, key, template.HTMLEscapeString(*value)))
		}
	}

	return fmt.Sprintf(`<input %s>`, strings.Join(attributes, " ")), nil
}

func (f *UUIDField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == nil {
		return "", nil
	}
	uuidValue, ok := value.(uuid.UUID)
	if !ok {
		return "", errors.New("value must be a uuid.UUID")
	}
	return form.HTMLType(uuidValue.String()), nil
}

func (f *UUIDField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	uuidValue, err := uuid.Parse(string(value))
	if err != nil {
		return nil, errors.New("invalid UUID format")
	}
	return uuidValue, nil
}

func (f *UUIDField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()

	baseValidations = append(baseValidations, f.requiredValidation)

	return baseValidations
}

func (f *UUIDField) requiredValidation(value interface{}) ([]error, error) {
	strValue, ok := value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	if f.Required && strValue == "" {
		return []error{errors.New("field is required")}, nil
	}
	return nil, nil
}

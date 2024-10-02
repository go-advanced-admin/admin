package fields

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"regexp"
	"strings"
)

type TextField struct {
	BaseField
	Placeholder *string
	MaxLength   *uint
	MinLength   *uint
	Required    bool
	Regex       *string
}

func (f *TextField) HTML() (string, error) {
	attributesMap := make(map[string]*string)
	if f.InitialValue != nil {
		htmlType, err := f.GoTypeToHTMLType(f.InitialValue)
		if err != nil {
			return "", err
		}
		value := template.HTMLEscapeString(string(htmlType))
		attributesMap["value"] = &value
	}
	if f.Placeholder != nil {
		value := template.HTMLEscapeString(*f.Placeholder)
		attributesMap["placeholder"] = &value
	}
	if f.MaxLength != nil {
		value := fmt.Sprintf("%d", *f.MaxLength)
		attributesMap["maxlength"] = &value
	}
	if f.MinLength != nil {
		value := fmt.Sprintf("%d", *f.MinLength)
		attributesMap["minlength"] = &value
	}
	if f.Required {
		attributesMap["required"] = nil
	}
	if f.Regex != nil {
		value := template.HTMLEscapeString(*f.Regex)
		attributesMap["pattern"] = &value
	}
	value := "text"
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
			attributes = append(attributes, fmt.Sprintf(`%s="%s""`, key, template.HTMLEscapeString(*value)))
		}
	}

	return fmt.Sprintf(`<input %s>`, strings.Join(attributes, " ")), nil
}

func (f *TextField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	strValue, ok := value.(string)
	if !ok {
		return "", errors.New("value must be a string")
	}
	return form.HTMLType(strValue), nil
}

func (f *TextField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	return string(value), nil
}

func (f *TextField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()

	baseValidations = append(baseValidations, f.requiredValidation, f.maxLengthValidation, f.minLengthValidation, f.regexValidation)

	return baseValidations
}

func (f *TextField) requiredValidation(value interface{}) ([]error, error) {
	strValue, ok := value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	if f.Required && strValue == "" {
		return []error{errors.New("field is required")}, nil
	}
	return nil, nil
}

func (f *TextField) maxLengthValidation(value interface{}) ([]error, error) {
	strValue, ok := value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	if f.MaxLength != nil && uint(len(strValue)) > *f.MaxLength {
		return []error{fmt.Errorf("input length of %d is greater than the maximum length of %d", len(strValue), *f.MaxLength)}, nil
	}
	return nil, nil
}

func (f *TextField) minLengthValidation(value interface{}) ([]error, error) {
	strValue, ok := value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	if f.MinLength != nil && uint(len(strValue)) < *f.MinLength {
		return []error{fmt.Errorf("input length of %d is less than the minimum length of %d", len(strValue), *f.MinLength)}, nil
	}
	return nil, nil
}

func (f *TextField) regexValidation(value interface{}) ([]error, error) {
	strValue, ok := value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	if f.Regex != nil {
		regex := *f.Regex
		matched, err := regexp.MatchString(regex, strValue)
		if err != nil {
			return nil, err
		}
		if !matched {
			return []error{errors.New("input does not match the required pattern")}, nil
		}
	}
	return nil, nil
}

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
	attributes := make([]string, 0)
	if f.InitialValue != nil {
		htmlType, err := f.GoTypeToHTMLType(f.InitialValue)
		if err != nil {
			return "", err
		}
		attributes = append(attributes, fmt.Sprintf(`value="%s"`, template.HTMLEscapeString(string(htmlType))))
	}
	if f.Placeholder != nil {
		attributes = append(attributes, fmt.Sprintf(`placeholder="%s"`, template.HTMLEscapeString(*f.Placeholder)))
	}
	if f.MaxLength != nil {
		attributes = append(attributes, fmt.Sprintf(`maxlength="%d"`, *f.MaxLength))
	}
	if f.MinLength != nil {
		attributes = append(attributes, fmt.Sprintf(`minlength="%d"`, *f.MinLength))
	}
	if f.Required {
		attributes = append(attributes, `required`)
	}
	if f.Regex != nil {
		attributes = append(attributes, fmt.Sprintf(`pattern="%s"`, template.HTMLEscapeString(*f.Regex)))
	}
	return fmt.Sprintf(`<input type="text" name="%s" %s>`, template.HTMLEscapeString(f.Name), strings.Join(attributes, " ")), nil
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

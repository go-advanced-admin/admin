package fields

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"strings"
	"time"
)

type DateField struct {
	BaseField
	Required    bool
	MinDate     *time.Time
	MaxDate     *time.Time
	Placeholder *string
}

func (f *DateField) HTML() (string, error) {
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

	if f.Required {
		attributesMap["required"] = nil
	}

	inputType := "date"
	attributesMap["type"] = &inputType
	name := template.HTMLEscapeString(f.Name)
	attributesMap["name"] = &name

	if f.MinDate != nil {
		value := f.MinDate.Format("2006-01-02")
		attributesMap["min"] = &value
	}

	if f.MaxDate != nil {
		value := f.MaxDate.Format("2006-01-02")
		attributesMap["max"] = &value
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

func (f *DateField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == nil {
		return "", nil
	}
	dateValue, ok := value.(time.Time)
	if !ok {
		return "", errors.New("value must be a time.Time")
	}
	return form.HTMLType(dateValue.Format("2006-01-02")), nil
}

func (f *DateField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	dateValue, err := time.Parse("2006-01-02", string(value))
	if err != nil {
		return nil, errors.New("invalid date format")
	}
	return dateValue, nil
}

func (f *DateField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()
	baseValidations = append(baseValidations, f.requiredValidation, f.minDateValidation, f.maxDateValidation)
	return baseValidations
}

func (f *DateField) requiredValidation(value interface{}) ([]error, error) {
	if !f.Required {
		return nil, nil
	}
	if value == nil {
		return []error{errors.New("field is required")}, nil
	}
	dateValue, ok := value.(time.Time)
	if !ok {
		return nil, errors.New("value must be a time.Time")
	}
	if dateValue.IsZero() {
		return []error{errors.New("field is required")}, nil
	}
	return nil, nil
}

func (f *DateField) minDateValidation(value interface{}) ([]error, error) {
	if value == nil {
		return nil, nil
	}
	dateValue, ok := value.(time.Time)
	if !ok {
		return nil, errors.New("value must be a time.Time")
	}
	if f.MinDate != nil && dateValue.Before(*f.MinDate) {
		return []error{fmt.Errorf("date %s is before minimum date %s", dateValue.Format("2006-01-02"), f.MinDate.Format("2006-01-02"))}, nil
	}
	return nil, nil
}

func (f *DateField) maxDateValidation(value interface{}) ([]error, error) {
	if value == nil {
		return nil, nil
	}
	dateValue, ok := value.(time.Time)
	if !ok {
		return nil, errors.New("value must be a time.Time")
	}
	if f.MaxDate != nil && dateValue.After(*f.MaxDate) {
		return []error{fmt.Errorf("date %s is after maximum date %s", dateValue.Format("2006-01-02"), f.MaxDate.Format("2006-01-02"))}, nil
	}
	return nil, nil
}

package fields

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"strconv"
	"strings"
)

type IntegerField struct {
	BaseField
	MinValue *int
	MaxValue *int
	Required bool
}

func (f *IntegerField) HTML() (string, error) {
	attributesMap := make(map[string]*string)

	if f.InitialValue != nil {
		intValue, err := f.GoTypeToHTMLType(f.InitialValue)
		if err != nil {
			return "", err
		}
		value := template.HTMLEscapeString(string(intValue))
		attributesMap["value"] = &value
	}

	if f.MinValue != nil {
		value := strconv.Itoa(*f.MinValue)
		attributesMap["min"] = &value
	}

	if f.MaxValue != nil {
		value := strconv.Itoa(*f.MaxValue)
		attributesMap["max"] = &value
	}

	if f.Required {
		attributesMap["required"] = nil
	}

	inputType := "number"
	attributesMap["type"] = &inputType
	name := template.HTMLEscapeString(f.Name)
	attributesMap["name"] = &name

	if f.SupersedingAttributes != nil {
		for k, v := range f.SupersedingAttributes {
			attributesMap[k] = v
		}
	}

	var attributes []string
	for key, value := range attributesMap {
		if value != nil {
			attributes = append(attributes, fmt.Sprintf(`%s="%s"`, key, template.HTMLEscapeString(*value)))
		} else {
			attributes = append(attributes, key)
		}
	}

	return fmt.Sprintf("<input %s>", strings.Join(attributes, " ")), nil
}

func (f *IntegerField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == "" {
		return "", nil
	}
	switch v := value.(type) {
	case int:
		return form.HTMLType(strconv.Itoa(v)), nil
	case int8:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	case int16:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	case int32:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	case int64:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	case uint:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	case uint8:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	case uint16:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	case uint32:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	case uint64:
		return form.HTMLType(strconv.Itoa(int(v))), nil
	default:
		return "", fmt.Errorf("value is not an int")
	}
}

func (f *IntegerField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	intValue, err := strconv.Atoi(string(value))
	if err != nil {
		return nil, err
	}
	return intValue, nil
}

func (f *IntegerField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()
	baseValidations = append(baseValidations, f.requiredValidation, f.minValidation, f.maxValidation)
	return baseValidations
}

func (f *IntegerField) requiredValidation(value interface{}) ([]error, error) {
	if f.Required && value == nil {
		return []error{fmt.Errorf("field %s is required", f.Name)}, nil
	}
	return nil, nil
}

func (f *IntegerField) minValidation(value interface{}) ([]error, error) {
	if f.MinValue != nil && value != nil {
		intValue, ok := value.(int)
		if !ok {
			return nil, fmt.Errorf("value is not an int")
		}
		if intValue < *f.MinValue {
			return []error{fmt.Errorf("field %s must be greater than %d", f.Name, *f.MinValue)}, nil
		}
	}
	return nil, nil
}

func (f *IntegerField) maxValidation(value interface{}) ([]error, error) {
	if f.MaxValue != nil && value != nil {
		intValue, ok := value.(int)
		if !ok {
			return nil, fmt.Errorf("value is not an int")
		}
		if intValue > *f.MaxValue {
			return []error{fmt.Errorf("field %s must be less than %d", f.Name, *f.MaxValue)}, nil
		}
	}
	return nil, nil
}

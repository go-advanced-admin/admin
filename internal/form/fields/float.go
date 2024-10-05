package fields

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"strconv"
	"strings"
)

type FloatField struct {
	BaseField
	MinValue *float64
	MaxValue *float64
	Required bool
}

func (f *FloatField) HTML() (string, error) {
	attributesMap := make(map[string]*string)

	if f.InitialValue != nil {
		floatValue, err := f.GoTypeToHTMLType(f.InitialValue)
		if err != nil {
			return "", err
		}
		value := template.HTMLEscapeString(string(floatValue))
		attributesMap["value"] = &value
	}

	if f.MinValue != nil {
		value := strconv.FormatFloat(*f.MinValue, 'f', -1, 64)
		attributesMap["min"] = &value
	}

	if f.MaxValue != nil {
		value := strconv.FormatFloat(*f.MaxValue, 'f', -1, 64)
		attributesMap["max"] = &value
	}

	if f.Required {
		attributesMap["required"] = nil
	}

	inputType := "number"
	attributesMap["type"] = &inputType
	step := "any"
	attributesMap["step"] = &step
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

func (f *FloatField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == nil {
		return "", nil
	}
	switch v := value.(type) {
	case float64:
		return form.HTMLType(strconv.FormatFloat(v, 'f', -1, 64)), nil
	case float32:
		return form.HTMLType(strconv.FormatFloat(float64(v), 'f', -1, 32)), nil
	default:
		return "", errors.New("value must be a float")
	}
}

func (f *FloatField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	floatValue, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return nil, errors.New("invalid float value")
	}
	return floatValue, nil
}

func (f *FloatField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()
	baseValidations = append(baseValidations, f.requiredValidation, f.minValueValidation, f.maxValueValidation)
	return baseValidations
}

func (f *FloatField) requiredValidation(value interface{}) ([]error, error) {
	if value == nil && f.Required {
		return []error{errors.New("field is required")}, nil
	}
	return nil, nil
}

func (f *FloatField) minValueValidation(value interface{}) ([]error, error) {
	if value != nil {
		floatValue, ok := value.(float64)
		if !ok {
			return nil, errors.New("value must be a float64")
		}
		if f.MinValue != nil && floatValue < *f.MinValue {
			return []error{fmt.Errorf("value %f is less than minimum %f", floatValue, *f.MinValue)}, nil
		}
	}
	return nil, nil
}

func (f *FloatField) maxValueValidation(value interface{}) ([]error, error) {
	if value != nil {
		floatValue, ok := value.(float64)
		if !ok {
			return nil, errors.New("value must be a float64")
		}
		if f.MaxValue != nil && floatValue > *f.MaxValue {
			return []error{fmt.Errorf("value %f is greater than maximum %f", floatValue, *f.MaxValue)}, nil
		}
	}
	return nil, nil
}

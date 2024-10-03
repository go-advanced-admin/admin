package fields

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"strings"
)

type MultipleChoiceField struct {
	BaseField
	Choices     []Choice
	Required    bool
	Placeholder *string
}

func (f *MultipleChoiceField) HTML() (string, error) {
	var options []string

	for _, choice := range f.Choices {
		selected := ""
		if f.InitialValue != nil {
			initialValues, ok := f.InitialValue.([]string)
			if !ok {
				return "", errors.New("initial value must be a slice of strings")
			}
			for _, val := range initialValues {
				if val == choice.Value {
					selected = " selected"
					break
				}
			}
		}
		option := fmt.Sprintf(`<option value="%s"%s>%s</option>`,
			template.HTMLEscapeString(choice.Value),
			selected,
			template.HTMLEscapeString(choice.Label),
		)
		options = append(options, option)
	}

	name := template.HTMLEscapeString(f.Name)
	html := fmt.Sprintf(`<select name="%s" multiple>%s</select>`, name, strings.Join(options, "\n"))
	return html, nil
}

func (f *MultipleChoiceField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == nil {
		return "", nil
	}
	values, ok := value.([]string)
	if !ok {
		return "", errors.New("value must be a slice of strings")
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return form.HTMLType(jsonValue), nil
}

func (f *MultipleChoiceField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" || value == "[]" {
		return nil, nil
	}
	var values []string
	if err := json.Unmarshal([]byte(value), &values); err != nil {
		return nil, errors.New("invalid multiple choice value")
	}
	return values, nil
}

func (f *MultipleChoiceField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()
	baseValidations = append(baseValidations, f.requiredValidation, f.choiceValidation)
	return baseValidations
}

func (f *MultipleChoiceField) requiredValidation(value interface{}) ([]error, error) {
	if !f.Required {
		return nil, nil
	}
	if value == nil {
		return []error{errors.New("at least one choice must be selected")}, nil
	}
	values, ok := value.([]string)
	if !ok {
		return nil, errors.New("value must be a slice of strings")
	}
	if f.Required && len(values) == 0 {
		return []error{errors.New("at least one choice must be selected")}, nil
	}
	return nil, nil
}

func (f *MultipleChoiceField) choiceValidation(value interface{}) ([]error, error) {
	if value == nil {
		return nil, nil
	}
	values, ok := value.([]string)
	if !ok {
		return nil, errors.New("value must be a slice of strings")
	}
	validChoices := make(map[string]bool)
	for _, choice := range f.Choices {
		validChoices[choice.Value] = true
	}
	for _, val := range values {
		if !validChoices[val] {
			return []error{fmt.Errorf("invalid choice: %s", val)}, nil
		}
	}
	return nil, nil
}

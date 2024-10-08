package fields

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"html/template"
	"strings"
)

type ChoiceField struct {
	BaseField
	Choices     []Choice
	Required    bool
	Placeholder *string
}

type Choice struct {
	Value string
	Label string
}

func (f *ChoiceField) HTML() (string, error) {
	var options []string
	if f.Placeholder != nil {
		options = append(options, fmt.Sprintf(`<option value="">%s</option>`, template.HTMLEscapeString(*f.Placeholder)))
	}

	for _, choice := range f.Choices {
		selected := ""
		if f.InitialValue != nil && f.InitialValue == choice.Value {
			selected = " selected"
		}
		option := fmt.Sprintf(`<option value="%s"%s>%s</option>`,
			template.HTMLEscapeString(choice.Value),
			selected,
			template.HTMLEscapeString(choice.Label),
		)
		options = append(options, option)
	}

	name := template.HTMLEscapeString(f.Name)
	html := fmt.Sprintf(`<select name="%s">%s</select>`, name, strings.Join(options, "\n"))
	return html, nil
}

func (f *ChoiceField) GoTypeToHTMLType(value interface{}) (form.HTMLType, error) {
	if value == nil {
		return "", nil
	}
	strValue, ok := value.(string)
	if !ok {
		return "", errors.New("value must be a string")
	}
	return form.HTMLType(strValue), nil
}

func (f *ChoiceField) HTMLTypeToGoType(value form.HTMLType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	return string(value), nil
}

func (f *ChoiceField) GetValidationFunctions() []form.FieldValidationFunc {
	baseValidations := f.BaseField.GetValidationFunctions()
	baseValidations = append(baseValidations, f.requiredValidation, f.choiceValidation)
	return baseValidations
}

func (f *ChoiceField) requiredValidation(value interface{}) ([]error, error) {
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
	if f.Required && strings.TrimSpace(strValue) == "" {
		return []error{errors.New("field is required")}, nil
	}
	return nil, nil
}

func (f *ChoiceField) choiceValidation(value interface{}) ([]error, error) {
	if value == nil {
		return nil, nil
	}
	strValue, ok := value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	for _, choice := range f.Choices {
		if strValue == choice.Value {
			return nil, nil
		}
	}
	return []error{errors.New("invalid choice selected")}, nil
}

package form

import (
	"fmt"
	"html/template"
	"strings"
)

type ValidationFunc func(map[string]interface{}) (frontend []error, backend error)

type Form interface {
	AddField(name string, field Field) error
	GetFields() []Field
	RegisterValidationFunctions(validationFuncs ...ValidationFunc)
	GetValidationFunctions() []ValidationFunc
	RegisterInitialValues(values map[string]interface{}) error
	Save(values map[string]HTMLType) (interface{}, error)
}

func ValuesAreValid(form Form, values map[string]interface{}) ([]error, map[string][]error, error) {
	formErrs := make([]error, 0)
	fieldsErrs := make(map[string][]error)

	fields := form.GetFields()
	for _, field := range fields {
		fieldName := field.GetName()
		fieldValue, exists := values[fieldName]
		if !exists {
			fieldValue = nil
		}
		fieldErrs, err := FieldValueIsValid(field, fieldValue)
		if err != nil {
			return formErrs, fieldsErrs, err
		}
		fieldsErrs[fieldName] = fieldErrs
	}

	validationFuncs := form.GetValidationFunctions()
	for _, validationFunc := range validationFuncs {
		frontend, err := validationFunc(values)
		if err != nil {
			return formErrs, fieldsErrs, err
		}
		if frontend != nil {
			formErrs = append(formErrs, frontend...)
		}
	}

	return formErrs, fieldsErrs, nil
}

func GetCleanData(form Form, values map[string]HTMLType) (map[string]interface{}, error) {
	cleanValues := make(map[string]interface{})

	for _, field := range form.GetFields() {
		fieldName := field.GetName()
		fieldValue, exists := values[fieldName]
		if !exists {
			fieldValue = ""
		}
		cleanValue, err := field.HTMLTypeToGoType(fieldValue)
		if err != nil {
			return nil, err
		}
		cleanValues[fieldName] = cleanValue
	}

	return cleanValues, nil
}

func RenderFormAsP(form Form) (string, error) {
	var htmlStrings []string
	for _, field := range form.GetFields() {
		fieldHTML, err := field.HTML()
		if err != nil {
			return "", err
		}
		label := template.HTMLEscapeString(field.GetLabel())
		htmlStrings = append(htmlStrings, fmt.Sprintf("<p><label>%s: %s</label></p>", label, fieldHTML))
	}
	return strings.Join(htmlStrings, "\n"), nil
}

func RenderFormAsUL(form Form) (string, error) {
	var htmlStrings []string
	for _, field := range form.GetFields() {
		fieldHTML, err := field.HTML()
		if err != nil {
			return "", err
		}
		label := template.HTMLEscapeString(field.GetLabel())
		htmlStrings = append(htmlStrings, fmt.Sprintf("<li><label>%s: %s</label></li>", label, fieldHTML))
	}
	return fmt.Sprintf("<ul>\n%s\n</ul>", strings.Join(htmlStrings, "\n")), nil
}

func RenderFormAsTable(form Form) (string, error) {
	var htmlStrings []string
	for _, field := range form.GetFields() {
		fieldHTML, err := field.HTML()
		if err != nil {
			return "", err
		}
		label := template.HTMLEscapeString(field.GetLabel())
		htmlStrings = append(htmlStrings, fmt.Sprintf("<tr><th><label>%s</label></th><td>%s</td></tr>", label, fieldHTML))
	}
	return fmt.Sprintf("<table>\n%s\n</table>", strings.Join(htmlStrings, "\n")), nil
}

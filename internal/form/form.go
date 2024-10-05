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

func renderErrors(errors []error) string {
	if len(errors) == 0 {
		return ""
	}
	var errStrings []string
	for _, err := range errors {
		errStrings = append(errStrings, template.HTMLEscapeString(err.Error()))
	}
	return fmt.Sprintf(`<ul class="errorlist"><li>%s</li></ul>`, strings.Join(errStrings, "</li><li>"))
}

func RenderFormAsP(form Form, formErrs []error, fieldsErrs map[string][]error) (string, error) {
	var htmlStrings []string
	for _, field := range form.GetFields() {
		fieldHTML, err := field.HTML()
		if err != nil {
			return "", err
		}
		label := template.HTMLEscapeString(field.GetLabel())
		fieldErrs, exists := fieldsErrs[field.GetName()]
		fieldErrors := ""
		if exists && len(fieldErrs) > 0 {
			fieldErrors = renderErrors(fieldErrs)
		}
		htmlStrings = append(htmlStrings, fmt.Sprintf(`<p><label for="%s">%s:</label> %s%s</p>`, field.GetName(), label, fieldHTML, fieldErrors))
	}
	if len(formErrs) > 0 {
		htmlStrings = append(htmlStrings, renderErrors(formErrs))
	}
	return strings.Join(htmlStrings, "\n"), nil
}

func RenderFormAsUL(form Form, formErrs []error, fieldsErrs map[string][]error) (string, error) {
	var htmlStrings []string
	for _, field := range form.GetFields() {
		fieldHTML, err := field.HTML()
		if err != nil {
			return "", err
		}
		label := template.HTMLEscapeString(field.GetLabel())
		fieldErrs, exists := fieldsErrs[field.GetName()]
		fieldErrors := ""
		if exists && len(fieldErrs) > 0 {
			fieldErrors = renderErrors(fieldErrs)
		}
		htmlStrings = append(htmlStrings, fmt.Sprintf(`<li><label for="%s">%s:</label> %s%s</li>`, field.GetName(), label, fieldHTML, fieldErrors))
	}
	if len(formErrs) > 0 {
		htmlStrings = append(htmlStrings, renderErrors(formErrs))
	}
	return fmt.Sprintf("<ul>\n%s\n</ul>", strings.Join(htmlStrings, "\n")), nil
}

func RenderFormAsTable(form Form, formErrs []error, fieldsErrs map[string][]error) (string, error) {
	var htmlStrings []string
	for _, field := range form.GetFields() {
		fieldHTML, err := field.HTML()
		if err != nil {
			return "", err
		}
		label := template.HTMLEscapeString(field.GetLabel())
		fieldErrs, exists := fieldsErrs[field.GetName()]
		fieldErrors := ""
		if exists && len(fieldErrs) > 0 {
			fieldErrors = renderErrors(fieldErrs)
		}
		htmlStrings = append(htmlStrings, fmt.Sprintf(`<tr><th><label for="%s">%s</label></th><td>%s%s</td></tr>`, field.GetName(), label, fieldHTML, fieldErrors))
	}
	if len(formErrs) > 0 {
		htmlStrings = append(htmlStrings, fmt.Sprintf("<tr><td colspan=\"2\">%s</td></tr>", renderErrors(formErrs)))
	}
	return fmt.Sprintf("<table>\n%s\n</table>", strings.Join(htmlStrings, "\n")), nil
}

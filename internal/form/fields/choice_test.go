package fields

import (
	"strings"
	"testing"
)

func TestChoiceField_HTML(t *testing.T) {
	placeholder := "Select an option"
	f := &ChoiceField{
		BaseField: BaseField{Name: "color"},
		Choices: []Choice{
			{Value: "red", Label: "Red"},
			{Value: "blue", Label: "Blue"},
		},
		Placeholder: &placeholder,
	}
	html, err := f.HTML()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(html, placeholder) {
		t.Error("Expected placeholder in HTML output")
	}

	f.InitialValue = "red"
	html, err = f.HTML()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(html, `value="red" selected`) {
		t.Error("Expected 'red' option to be selected")
	}
}

func TestChoiceField_GoTypeToHTMLType(t *testing.T) {
	f := &ChoiceField{}
	htmlType, err := f.GoTypeToHTMLType("value")
	if err != nil || htmlType != "value" {
		t.Errorf("Expected 'value', got '%s' with error '%v'", htmlType, err)
	}

	htmlType, err = f.GoTypeToHTMLType(nil)
	if err != nil || htmlType != "" {
		t.Errorf("Expected empty string, got '%s' with error '%v'", htmlType, err)
	}

	_, err = f.GoTypeToHTMLType(123)
	if err == nil {
		t.Error("Expected error when value is not a string")
	}
}

func TestChoiceField_HTMLTypeToGoType(t *testing.T) {
	f := &ChoiceField{}
	val, err := f.HTMLTypeToGoType("value")
	if err != nil || val != "value" {
		t.Errorf("Expected 'value', got '%v' with error '%v'", val, err)
	}

	val, err = f.HTMLTypeToGoType("")
	if err != nil || val != nil {
		t.Errorf("Expected nil, got '%v' with error '%v'", val, err)
	}
}

func TestChoiceField_requiredValidation(t *testing.T) {
	f := &ChoiceField{Required: true}
	errs, err := f.requiredValidation(nil)
	if err != nil {
		t.Errorf("Unexpected backend error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("Expected frontend error for required field")
	}

	f.Required = false
	errs, err = f.requiredValidation("")
	if err != nil || len(errs) != 0 {
		t.Errorf("Expected no errors, got errs: %v, backend error: %v", errs, err)
	}
}

func TestChoiceField_choiceValidation(t *testing.T) {
	f := &ChoiceField{
		Choices: []Choice{
			{Value: "red", Label: "Red"},
			{Value: "blue", Label: "Blue"},
		},
	}
	errs, err := f.choiceValidation("red")
	if err != nil || len(errs) != 0 {
		t.Errorf("Expected no errors, got errs: %v, backend error: %v", errs, err)
	}

	errs, err = f.choiceValidation("green")
	if err != nil {
		t.Errorf("Unexpected backend error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("Expected frontend error for invalid choice")
	}
}

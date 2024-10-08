package fields

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBooleanField_HTML(t *testing.T) {
	f := &BooleanField{
		BaseField: BaseField{Name: "agree", InitialValue: true},
		Required:  true,
	}
	html, err := f.HTML()
	assert.Nil(t, err)
	assert.Contains(t, html, "type=\"checkbox\"")
	assert.Contains(t, html, "checked")
	assert.Contains(t, html, "required")

	f.InitialValue = "not a bool"
	_, err = f.HTML()
	assert.Error(t, err)
}

func TestBooleanField_GoTypeToHTMLType(t *testing.T) {
	f := &BooleanField{}
	htmlType, err := f.GoTypeToHTMLType(true)
	if err != nil || htmlType != "on" {
		t.Errorf("Expected 'on', got '%s' with error '%v'", htmlType, err)
	}

	htmlType, err = f.GoTypeToHTMLType(false)
	if err != nil || htmlType != "" {
		t.Errorf("Expected '', got '%s' with error '%v'", htmlType, err)
	}

	_, err = f.GoTypeToHTMLType("not a bool")
	if err == nil {
		t.Error("Expected error when value is not a boolean")
	}
}

func TestBooleanField_HTMLTypeToGoType(t *testing.T) {
	f := &BooleanField{}
	val, err := f.HTMLTypeToGoType("on")
	if err != nil || val != true {
		t.Errorf("Expected true, got '%v' with error '%v'", val, err)
	}

	val, err = f.HTMLTypeToGoType("")
	if err != nil || val != nil {
		t.Errorf("Expected nil, got '%v' with error '%v'", val, err)
	}

	val, err = f.HTMLTypeToGoType("false")
	if err != nil || val != false {
		t.Errorf("Expected false, got '%v' with error '%v'", val, err)
	}
}

func TestBooleanField_GetValidationFunctions(t *testing.T) {
	f := &BooleanField{}
	funcs := f.GetValidationFunctions()
	if len(funcs) != 1 {
		t.Errorf("Expected 1 validation function, got %d", len(funcs))
	}
}

func TestBooleanField_requiredValidation(t *testing.T) {
	f := &BooleanField{Required: true}
	errs, err := f.requiredValidation(nil)
	if err != nil {
		t.Errorf("Unexpected backend error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("Expected frontend error for required field")
	}

	f.Required = false
	errs, err = f.requiredValidation(nil)
	if err != nil || len(errs) != 0 {
		t.Errorf("Expected no errors, got errs: %v, backend error: %v", errs, err)
	}
}

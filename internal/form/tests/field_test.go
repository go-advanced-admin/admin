package tests

import (
	"errors"
	"github.com/go-advanced-admin/admin/internal/form"
	"github.com/go-advanced-admin/admin/internal/form/fields"
	"github.com/stretchr/testify/assert"
	"testing"
)

func DummyFieldValidationFunc(value interface{}) ([]error, error) {
	strValue, ok := value.(string)
	if !ok || strValue == "" {
		return []error{errors.New("invalid field value")}, nil
	}
	return nil, nil
}

func TestFieldValueIsValid(t *testing.T) {
	textField := &fields.TextField{}
	err := textField.RegisterName("exampleField")
	assert.Nil(t, err)
	textField.RegisterValidationFunctions(DummyFieldValidationFunc)

	fieldErrs, err := form.FieldValueIsValid(textField, "")
	assert.NotNil(t, fieldErrs)
	assert.Nil(t, err)
	assert.Contains(t, fieldErrs, errors.New("invalid field value"))

	fieldErrs, err = form.FieldValueIsValid(textField, "valid")
	assert.Nil(t, fieldErrs)
	assert.Nil(t, err)
}

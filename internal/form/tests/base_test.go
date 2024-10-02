package tests

import (
	"errors"
	"github.com/go-advanced-admin/admin/internal/form/fields"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseFieldMethods(t *testing.T) {
	baseField := &fields.BaseField{}
	err := baseField.RegisterName("example")
	assert.Nil(t, err)
	assert.Equal(t, "example", baseField.GetName())

	err = baseField.RegisterName("")
	assert.NotNil(t, err)

	baseField.RegisterInitialValue("initial")
	assert.Equal(t, "initial", baseField.InitialValue)

	validateFunc := func(value interface{}) ([]error, error) {
		if value == "invalid" {
			return []error{errors.New("value is invalid")}, nil
		}
		return nil, nil
	}

	baseField.RegisterValidationFunctions(validateFunc)
	assert.NotNil(t, baseField.GetValidationFunctions())
}

package fields

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseField_GetName(t *testing.T) {
	f := &BaseField{
		Name: "test_name",
	}
	assert.Equal(t, "test_name", f.GetName())
}

func TestBaseField_RegisterInitialValue(t *testing.T) {
	f := &BaseField{}
	f.RegisterInitialValue(42)
	assert.Equal(t, 42, f.InitialValue)
}

func TestBaseField_GetValidationFunctions(t *testing.T) {
	f := &BaseField{}
	funcs := f.GetValidationFunctions()
	assert.NotNil(t, funcs)
	assert.Len(t, funcs, 0)

	validationFunc := func(value interface{}) ([]error, error) {
		return nil, nil
	}
	f.RegisterValidationFunctions(validationFunc)
	funcs = f.GetValidationFunctions()
	assert.NotNil(t, funcs)
	assert.Len(t, funcs, 1)
}

func TestBaseField_RegisterName(t *testing.T) {
	f := &BaseField{}
	err := f.RegisterName("test_name")
	assert.Nil(t, err)
	assert.Equal(t, "test_name", f.Name)

	err = f.RegisterName("")
	assert.Error(t, err)
}

func TestBaseField_RegisterLabel(t *testing.T) {
	f := &BaseField{}
	err := f.RegisterLabel("Test Label")
	assert.Nil(t, err)
	assert.Equal(t, "Test Label", f.Label)

	err = f.RegisterLabel("")
	assert.Error(t, err)
}

func TestBaseField_GetLabel(t *testing.T) {
	f := &BaseField{Name: "testName", Label: "Test Label"}
	assert.Equal(t, "Test Label", f.GetLabel())

	f.Label = ""
	assert.Equal(t, "Test Name", f.GetLabel())
}

func TestBaseField_RegisterValidationFunctions(t *testing.T) {
	f := &BaseField{}
	validationFunc1 := func(value interface{}) ([]error, error) { return nil, nil }
	validationFunc2 := func(value interface{}) ([]error, error) { return nil, nil }
	f.RegisterValidationFunctions(validationFunc1, validationFunc2)
	assert.Len(t, f.ValidationFuncs, 2)
}

func TestBaseField_SetSupersedingAttribute(t *testing.T) {
	f := &BaseField{}
	value1 := "value1"
	f.SetSupersedingAttribute("attr1", &value1)
	assert.NotNil(t, f.SupersedingAttributes)
	assert.Len(t, f.SupersedingAttributes, 1)
	assert.Equal(t, &value1, f.SupersedingAttributes["attr1"])
	assert.Equal(t, "value1", *f.SupersedingAttributes["attr1"])

	value2 := "value2"
	f.SetSupersedingAttribute("attr1", &value2)
	assert.NotNil(t, f.SupersedingAttributes)
	assert.Len(t, f.SupersedingAttributes, 1)
	assert.Equal(t, &value2, f.SupersedingAttributes["attr1"])
	assert.Equal(t, "value2", *f.SupersedingAttributes["attr1"])

	value3 := "value3"
	f.SetSupersedingAttribute("attr3", &value3)
	assert.NotNil(t, f.SupersedingAttributes)
	assert.Len(t, f.SupersedingAttributes, 2)
	assert.Equal(t, &value2, f.SupersedingAttributes["attr1"])
	assert.Equal(t, "value2", *f.SupersedingAttributes["attr1"])
	assert.Equal(t, &value3, f.SupersedingAttributes["attr3"])
	assert.Equal(t, "value3", *f.SupersedingAttributes["attr3"])
}

func TestBaseFieldMethods(t *testing.T) {
	baseField := &BaseField{}
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

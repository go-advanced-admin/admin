package tests

import (
	"errors"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"github.com/go-advanced-admin/admin/internal/form/fields"
	"github.com/go-advanced-admin/admin/internal/form/forms"
	"github.com/stretchr/testify/assert"
	"testing"
)

func dummyValidationFunc(values map[string]interface{}) ([]error, error) {
	if name, exists := values["name"].(string); exists && name == "invalid" {
		return []error{errors.New("invalid name")}, nil
	}
	return nil, nil
}

func TestValuesAreValid(t *testing.T) {
	textField1 := &fields.TextField{}
	textField2 := &fields.TextField{}

	baseForm := &forms.BaseForm{}
	err := baseForm.AddField("name", textField1)
	assert.Nil(t, err)
	err = baseForm.AddField("age", textField2)
	assert.Nil(t, err)
	baseForm.RegisterValidationFunctions(dummyValidationFunc)

	values := map[string]interface{}{
		"name": "invalid",
		"age":  "25",
	}

	formErrs, fieldErrs, err := form.ValuesAreValid(baseForm, values)
	assert.NotNil(t, formErrs)
	assert.Nil(t, err)
	assert.NotEmpty(t, fieldErrs)
	assert.Empty(t, fieldErrs["name"])
	assert.Empty(t, fieldErrs["age"])
	assert.Contains(t, fmt.Sprint(formErrs), "invalid name")
}

func TestGetCleanData(t *testing.T) {
	textField1 := &fields.TextField{}
	textField2 := &fields.TextField{}

	baseForm := &forms.BaseForm{}
	err := baseForm.AddField("name", textField1)
	assert.Empty(t, err)
	err = baseForm.AddField("age", textField2)
	assert.Empty(t, err)

	values := map[string]form.HTMLType{
		"name":  "John Doe",
		"age":   "30",
		"extra": "thisshouldfail",
	}

	cleanValues, err := form.GetCleanData(baseForm, values)
	assert.Nil(t, err)
	assert.NotNil(t, cleanValues)
	_, exists := cleanValues["extra"]
	assert.False(t, exists)
	assert.Equal(t, "John Doe", cleanValues["name"])
	assert.Equal(t, "30", cleanValues["age"])

	values = map[string]form.HTMLType{
		"name": "Jane Doe",
		"age":  "29",
	}

	cleanValues, err = form.GetCleanData(baseForm, values)
	assert.Nil(t, err)
	assert.NotNil(t, cleanValues)
	assert.Equal(t, "Jane Doe", cleanValues["name"])
	assert.Equal(t, "29", cleanValues["age"])
}

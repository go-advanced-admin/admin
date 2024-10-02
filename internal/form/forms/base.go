package forms

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"log"
)

type BaseForm struct {
	Fields          []form.Field
	ValidationFuncs []form.ValidationFunc
}

func (b *BaseForm) GetField(name string) (form.Field, bool) {
	for _, field := range b.Fields {
		if field.GetName() == name {
			return field, true
		}
	}
	return nil, false
}

func (b *BaseForm) AddField(name string, field form.Field) error {
	err := field.RegisterName(name)
	if err != nil {
		return err
	}
	_, exists := b.GetField(name)
	if exists {
		return fmt.Errorf("field %s already exists", name)
	}
	b.Fields = append(b.Fields, field)
	return nil
}

func (b *BaseForm) GetFields() []form.Field {
	return b.Fields
}

func (b *BaseForm) RegisterValidationFunctions(validationFuncs ...form.ValidationFunc) {
	b.ValidationFuncs = append(b.ValidationFuncs, validationFuncs...)
}

func (b *BaseForm) GetValidationFunctions() []form.ValidationFunc {
	return b.ValidationFuncs
}

// RegisterInitialValues TODO make it return an error if there are values in the map that do not match any field
func (b *BaseForm) RegisterInitialValues(values map[string]interface{}) error {
	for _, field := range b.Fields {
		if value, exists := values[field.GetName()]; exists {
			field.RegisterInitialValue(value)
		}
	}
	return nil
}

func (b *BaseForm) Save(values map[string]form.HTMLType) (string, error) {
	log.Println("Saving form", values)
	cleanValues, err := form.GetCleanData(b, values)
	if err != nil {
		return "", err
	}
	log.Println("Clean values", cleanValues)
	return "", nil
}

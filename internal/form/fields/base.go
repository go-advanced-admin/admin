package fields

import (
	"errors"
	"github.com/go-advanced-admin/admin/internal/form"
)

type BaseField struct {
	Name                  string
	Label                 string
	ValidationFuncs       []form.FieldValidationFunc
	InitialValue          interface{}
	SupersedingAttributes map[string]*string
}

func (f *BaseField) GetName() string {
	return f.Name
}

func (f *BaseField) RegisterInitialValue(value interface{}) {
	f.InitialValue = value
}

func (f *BaseField) GetValidationFunctions() []form.FieldValidationFunc {
	if f.ValidationFuncs == nil {
		return make([]form.FieldValidationFunc, 0)
	}
	return f.ValidationFuncs
}

func (f *BaseField) RegisterName(name string) error {
	if name == "" {
		return errors.New("field name cannot be empty")
	}
	f.Name = name
	return nil
}

func (f *BaseField) RegisterLabel(label string) error {
	if label == "" {
		return errors.New("field label cannot be empty")
	}
	f.Label = label
	return nil
}

func (f *BaseField) GetLabel() string {
	return f.Label
}

func (f *BaseField) RegisterValidationFunctions(validationFuncs ...form.FieldValidationFunc) {
	if f.ValidationFuncs == nil {
		f.ValidationFuncs = make([]form.FieldValidationFunc, 0)
	}
	f.ValidationFuncs = append(f.ValidationFuncs, validationFuncs...)
}

func (f *BaseField) SetSupersedingAttribute(name string, value *string) {
	if f.SupersedingAttributes == nil {
		f.SupersedingAttributes = make(map[string]*string)
	}
	f.SupersedingAttributes[name] = value
}

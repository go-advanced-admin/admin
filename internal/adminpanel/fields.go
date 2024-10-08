package adminpanel

import (
	"github.com/go-advanced-admin/admin/internal/form"
	"reflect"
)

// FieldConfig holds configuration for a model field in the admin panel.
type FieldConfig struct {
	Name                  string
	DisplayName           string
	FieldType             reflect.Type
	IncludeInListFetch    bool
	IncludeInListDisplay  bool
	IncludeInSearch       bool
	IncludeInInstanceView bool
	AddFormField          form.Field
	EditFormField         form.Field
}

// AdminFormFieldInterface allows a model to customize form fields for add and edit operations.
type AdminFormFieldInterface interface {
	// AdminFormField returns a custom form field for the given field name and operation (isEdit).
	AdminFormField(name string, isEdit bool) form.Field
}

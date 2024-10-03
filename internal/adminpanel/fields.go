package adminpanel

import (
	"github.com/go-advanced-admin/admin/internal/form"
	"reflect"
)

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

type AdminFormFieldInterface interface {
	AdminFormField(name string, isEdit bool) form.Field
}

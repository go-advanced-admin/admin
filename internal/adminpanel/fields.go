package adminpanel

import "reflect"

type FieldConfig struct {
	Name                  string
	DisplayName           string
	FieldType             reflect.Type
	IncludeInListFetch    bool
	IncludeInListDisplay  bool
	IncludeInSearch       bool
	IncludeInInstanceView bool
	IncludeInAddForm      bool
	IncludeInEditForm     bool
}

package utils

import (
	"fmt"
	"reflect"
)

func GetFieldValue(instance interface{}, fieldName string) (interface{}, error) {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fieldValue := val.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	return fieldValue.Interface(), nil
}

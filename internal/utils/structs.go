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

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("instance is not a struct or pointer to a struct")
	}

	fieldValue := val.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	return fieldValue.Interface(), nil
}

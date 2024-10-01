package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func SetStringsAsType(value reflect.Value, input string) error {
	switch value.Kind() {
	case reflect.String:
		value.SetString(input)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return err
		}
		value.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return err
		}
		value.SetUint(parsed)
	default:
		return fmt.Errorf("unsupported type: %s", value.Type())
	}

	return nil
}

package utils

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strconv"
)

func SetStringsAsType(value reflect.Value, input string) error {
	if !value.CanSet() {
		return fmt.Errorf("value is not settable")
	}

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
		if value.Type() == reflect.TypeOf(uuid.UUID{}) {
			u, err := uuid.Parse(input)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(u))
			return nil
		}
		return fmt.Errorf("unsupported type: %s", value.Type())
	}

	return nil
}

func ConvertStringToType(s string, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.String:
		return s, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, t.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(i).Convert(t).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, t.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(u).Convert(t).Interface(), nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(f).Convert(t).Interface(), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		if t == reflect.TypeOf(uuid.UUID{}) {
			u, err := uuid.Parse(s)
			if err != nil {
				return nil, err
			}
			return u, nil
		}
		return nil, errors.New("unsupported kind: " + t.Kind().String())
	}
}

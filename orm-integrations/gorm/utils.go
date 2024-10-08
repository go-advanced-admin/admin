package admingorm

import (
	"reflect"
	"strings"
)

func getGormColumnName(structField reflect.StructField) string {
	tag, ok := structField.Tag.Lookup("gorm")
	if !ok {
		return structField.Name
	}
	tagParts := strings.Split(tag, ";")
	for _, part := range tagParts {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return structField.Name
}

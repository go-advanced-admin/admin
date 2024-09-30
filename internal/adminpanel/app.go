package adminpanel

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type App struct {
	Name   string
	Models map[string]*Model
}

func (a *App) RegisterModel(model interface{}) (*Model, error) {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("admin model '%s' must be a pointer to a struct", modelType.Name())
	}

	modelType = modelType.Elem()
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("admin model '%s' must be a pointer to a struct", modelType.Name())
	}

	var name string
	namer, ok := model.(AdminModelNameInterface)
	if ok {
		name = namer.AdminName()
	} else {
		name = humanizeName(modelType.Name())
	}

	if _, exists := a.Models[name]; exists {
		return nil, fmt.Errorf("admin model '%s' already exists in app '%s'. Models cannot be registered more than once", name, a.Name)
	}
	a.Models[name] = &Model{Name: name, PTR: model}
	return a.Models[name], nil
}

func humanizeName(name string) string {
	var result []rune
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) && !(unicode.IsUpper(r) && unicode.IsUpper(rune(name[i-1]))) {
			result = append(result, ' ')
		}
		result = append(result, r)
	}
	return strings.TrimSpace(string(result))
}

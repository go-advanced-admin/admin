package adminpanel

import (
	"fmt"
	"net/http"
	"reflect"
)

type FieldConfig struct {
	Name                 string
	DisplayName          string
	IncludeInListFetch   bool
	IncludeInListDisplay bool
}

type Model struct {
	Name             string
	DisplayName      string
	PTR              interface{}
	App              *App
	Fields           []FieldConfig
	PrimaryKeyGetter func(interface{}) interface{}
}

type AdminModelNameInterface interface {
	AdminName() string
}

type AdminModelDisplayNameInterface interface {
	AdminDisplayName() string
}

type AdminModelGetIDInterface interface {
	AdminGetID() interface{}
}

func (m *Model) GetLink() string {
	return fmt.Sprintf("%s/%s", m.App.GetLink(), m.Name)
}

func (m *Model) GetFullLink() string {
	return m.App.Panel.Config.GetLink(m.GetLink())
}

func (m *Model) GetViewHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		allowed, err := m.App.Panel.PermissionChecker.HasModelReadPermission(m.App.Name, m.Name, data)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		if !allowed {
			return http.StatusForbidden, "Forbidden"
		}

		apps, err := GetAppsWithReadPermissions(m.App.Panel, data)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		var fieldsToFetch []string
		for _, fieldConfig := range m.Fields {
			if fieldConfig.IncludeInListFetch {
				fieldsToFetch = append(fieldsToFetch, fieldConfig.Name)
			}
		}

		instances, err := m.App.Panel.ORM.FetchInstancesOnlyFields(m.PTR, fieldsToFetch)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		filteredInstances, err := filterInstancesByPermission(instances, m, data)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		html, err := m.App.Panel.Config.Renderer.RenderTemplate("model.html", map[string]interface{}{"apps": apps, "model": m, "instances": filteredInstances})
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusOK, html
	}
}

func GetPrimaryKeyGetter(model interface{}) (func(interface{}) interface{}, error) {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	if _, implements := model.(AdminModelGetIDInterface); implements {
		return func(instance interface{}) interface{} {
			return instance.(AdminModelGetIDInterface).AdminGetID()
		}, nil
	}

	if idField, found := modelType.Elem().FieldByName("ID"); found {
		return func(instance interface{}) interface{} {
			return modelValue.Elem().FieldByName(idField.Name).Interface()
		}, nil
	}

	return nil, fmt.Errorf("no valid primary key method or ID field found. A struct must either have the ID field or implement func AdminGetID() interface{}")
}

func filterInstancesByPermission(instances interface{}, model *Model, data interface{}) ([]interface{}, error) {
	val := reflect.ValueOf(instances)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	filtered := make([]interface{}, 0, val.Len())

	for i := 0; i < val.Len(); i++ {
		instance := val.Index(i).Interface()
		id := model.PrimaryKeyGetter(instance)
		allowed, err := model.App.Panel.PermissionChecker.HasInstanceReadPermission(model.App.Name, model.Name, id, data)
		if err != nil {
			return nil, err
		}
		if allowed && instance != nil {
			filtered = append(filtered, instance)
		}
	}

	return filtered, nil
}

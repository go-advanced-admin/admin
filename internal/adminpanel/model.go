package adminpanel

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

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

type Model struct {
	Name             string
	DisplayName      string
	PTR              interface{}
	App              *App
	Fields           []FieldConfig
	PrimaryKeyGetter func(interface{}) interface{}
	PrimaryKeyType   reflect.Type
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

func (m *Model) GetAddLink() string {
	return fmt.Sprintf("%s/add", m.GetLink())
}

func (m *Model) GetFullAddLink() string {
	return m.App.Panel.Config.GetLink(m.GetAddLink())
}

func (m *Model) GetViewHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		var page, perPage uint
		pageQuery := m.App.Panel.Web.GetQueryParam(data, "page")
		perPageQuery := m.App.Panel.Web.GetQueryParam(data, "perPage")

		if p, err := strconv.Atoi(pageQuery); err == nil {
			page = uint(p)
		} else {
			page = 1
		}

		if pp, err := strconv.Atoi(perPageQuery); err == nil {
			perPage = uint(pp)
		} else {
			perPage = m.App.Panel.Config.DefaultInstancesPerPage
		}

		if perPage < 10 {
			perPage = 10
		}

		allowed, err := m.App.Panel.PermissionChecker.HasModelReadPermission(m.App.Name, m.Name, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("forbidden"))
		}

		apps, err := GetAppsWithReadPermissions(m.App.Panel, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		var fieldsToFetch []string
		for _, fieldConfig := range m.Fields {
			if fieldConfig.IncludeInListFetch {
				fieldsToFetch = append(fieldsToFetch, fieldConfig.Name)
			}
		}

		searchQuery := m.App.Panel.Web.GetQueryParam(data, "search")
		var instances interface{}
		if searchQuery == "" {
			instances, err = m.App.Panel.ORM.FetchInstancesOnlyFields(m.PTR, fieldsToFetch)
		} else {
			var fieldsToSearch []string
			for _, fieldConfig := range m.Fields {
				if fieldConfig.IncludeInSearch {
					fieldsToSearch = append(fieldsToSearch, fieldConfig.Name)
				}
			}
			instances, err = m.App.Panel.ORM.FetchInstancesOnlyFieldWithSearch(m.PTR, fieldsToFetch, searchQuery, fieldsToSearch)
		}
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		filteredInstances, err := filterInstancesByPermission(instances, m, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		totalCount := uint(len(filteredInstances))
		totalPages := (totalCount + perPage - 1) / perPage

		startIndex := (page - 1) * perPage
		endIndex := startIndex + perPage

		if startIndex > totalCount {
			startIndex = totalCount
		}
		if endIndex > totalCount {
			endIndex = totalCount
		}

		pagedInstances := filteredInstances[startIndex:endIndex]

		cleanInstances := make([]Instance, len(pagedInstances))
		for i, instance := range pagedInstances {
			id := m.PrimaryKeyGetter(instance)
			updateAllowed, err := m.App.Panel.PermissionChecker.HasInstanceUpdatePermission(m.App.Name, m.Name, id, data)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			deleteAllowed, err := m.App.Panel.PermissionChecker.HasInstanceDeletePermission(m.App.Name, m.Name, id, data)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			cleanInstance := Instance{
				InstanceID:  id,
				Data:        instance,
				Model:       m,
				Permissions: Permissions{Read: true, Update: updateAllowed, Delete: deleteAllowed},
			}
			cleanInstances[i] = cleanInstance
		}

		html, err := m.App.Panel.Config.Renderer.RenderTemplate("model", map[string]interface{}{
			"apps":        apps,
			"model":       m,
			"instances":   cleanInstances,
			"totalCount":  totalCount,
			"totalPages":  totalPages,
			"currentPage": page,
			"perPage":     perPage,
		})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		return http.StatusOK, html
	}
}

func GetPrimaryKeyGetter(model interface{}) (func(interface{}) interface{}, error) {
	modelType := reflect.TypeOf(model)

	if _, implements := model.(AdminModelGetIDInterface); implements {
		return func(instance interface{}) interface{} {
			return instance.(AdminModelGetIDInterface).AdminGetID()
		}, nil
	}

	if idField, found := modelType.Elem().FieldByName("ID"); found {
		return func(instance interface{}) interface{} {
			instValue := reflect.ValueOf(instance)

			if instValue.Kind() == reflect.Ptr {
				instValue = instValue.Elem()
			}

			if !instValue.FieldByName(idField.Name).IsValid() {
				panic("ID field does not exist in instance")
			}

			return instValue.FieldByName(idField.Name).Interface()
		}, nil
	}

	return nil, fmt.Errorf("no valid primary key method or ID field found. A struct must either have the ID field or implement func AdminGetID() interface{}")
}

func filterInstancesByPermission(instances interface{}, model *Model, data interface{}) ([]interface{}, error) {
	val := reflect.ValueOf(instances)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil, fmt.Errorf("instances must be a slice or array")
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

package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/logging"
	"net/http"
	"reflect"
	"strconv"
)

// Model represents a registered model within an app in the admin panel.
type Model struct {
	Name        string
	DisplayName string
	PTR         interface{}
	App         *App
	Fields      []FieldConfig
	ORM         ORMIntegrator
}

// CreateViewLog creates a log entry when the model's list view is accessed.
func (m *Model) CreateViewLog(ctx interface{}) error {
	return m.App.Panel.Config.CreateLog(ctx, logging.LogStoreLevelListView, fmt.Sprintf("%s | %s", m.App.Name, m.DisplayName), nil, "", "")
}

// GetORM returns the ORM integrator for the model.
func (m *Model) GetORM() ORMIntegrator {
	if m.ORM != nil {
		return m.ORM
	}
	return m.App.GetORM()
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

// GetLink returns the relative URL path to the model.
func (m *Model) GetLink() string {
	return fmt.Sprintf("%s/%s", m.App.GetLink(), m.Name)
}

// GetFullLink returns the full URL path to the model, including the admin prefix.
func (m *Model) GetFullLink() string {
	return m.App.Panel.Config.GetLink(m.GetLink())
}

// GetAddLink returns the relative URL path to add a new instance of the model.
func (m *Model) GetAddLink() string {
	return fmt.Sprintf("%s/add", m.GetLink())
}

// GetFullAddLink returns the full URL path to add a new instance of the model.
func (m *Model) GetFullAddLink() string {
	return m.App.Panel.Config.GetLink(m.GetAddLink())
}

// GetViewHandler returns the HTTP handler function for the model's list view.
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
			instances, err = m.GetORM().FetchInstancesOnlyFields(m.PTR, fieldsToFetch)
		} else {
			var fieldsToSearch []string
			for _, fieldConfig := range m.Fields {
				if fieldConfig.IncludeInSearch {
					fieldsToSearch = append(fieldsToSearch, fieldConfig.Name)
				}
			}
			instances, err = m.GetORM().FetchInstancesOnlyFieldWithSearch(m.PTR, fieldsToFetch, searchQuery, fieldsToSearch)
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
			id, err := m.GetPrimaryKeyValue(instance)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
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
			"navBarItems": m.App.Panel.Config.GetNavBarItems(data),
		})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		err = m.CreateViewLog(data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		return http.StatusOK, html
	}
}

// GetPrimaryKeyValue retrieves the primary key value of an instance.
func (m *Model) GetPrimaryKeyValue(instance interface{}) (interface{}, error) {
	return m.GetORM().GetPrimaryKeyValue(instance)
}

// GetPrimaryKeyType retrieves the primary key type of the model.
func (m *Model) GetPrimaryKeyType() (reflect.Type, error) {
	return m.GetORM().GetPrimaryKeyType(m.PTR)
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
		id, err := model.GetPrimaryKeyValue(instance)
		if err != nil {
			return nil, err
		}
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

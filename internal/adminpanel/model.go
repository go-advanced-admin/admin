package adminpanel

import (
	"fmt"
	"log"
	"net/http"
)

type FieldConfig struct {
	Name                 string
	DisplayName          string
	IncludeInListDisplay bool
}

type Model struct {
	Name        string
	DisplayName string
	PTR         interface{}
	App         *App
	Fields      []FieldConfig
}

type AdminModelNameInterface interface {
	AdminName() string
}

type AdminModelDisplayNameInterface interface {
	AdminDisplayName() string
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
			if fieldConfig.IncludeInListDisplay {
				fieldsToFetch = append(fieldsToFetch, fieldConfig.Name)
			}
		}

		instances, err := m.App.Panel.ORM.FetchInstancesOnlyFields(m.PTR, fieldsToFetch)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		log.Println(instances)

		html, err := m.App.Panel.Config.Renderer.RenderTemplate("model.html", map[string]interface{}{"apps": apps, "model": m, "instances": instances})
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusOK, html
	}
}

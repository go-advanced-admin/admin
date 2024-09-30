package adminpanel

import (
	"fmt"
	"net/http"
)

type Model struct {
	Name        string
	DisplayName string
	PTR         interface{}
	App         *App
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

		html, err := m.App.Panel.Config.Renderer.RenderTemplate("model.html", map[string]interface{}{"apps": apps, "model": m})
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusOK, html
	}
}

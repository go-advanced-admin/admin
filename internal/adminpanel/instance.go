package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/utils"
	"net/http"
	"reflect"
)

type Instance struct {
	InstanceID  interface{}
	Data        interface{}
	Model       *Model
	Permissions Permissions
}

func (i *Instance) GetLink() string {
	return fmt.Sprintf("%s/%v", i.Model.GetLink(), i.InstanceID)
}

func (i *Instance) GetFullLink() string {
	return i.Model.App.Panel.Config.GetLink(i.GetLink())
}

func (i *Instance) GetEditLink() string {
	return fmt.Sprintf("%s/edit", i.GetLink())
}

func (i *Instance) GetFullEditLink() string {
	return i.Model.App.Panel.Config.GetLink(i.GetEditLink())
}

func (m *Model) GetInstanceDeleteHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := m.App.Panel.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		primaryKeyValue := reflect.New(m.PrimaryKeyType).Elem()
		if err := utils.SetStringsAsType(primaryKeyValue, instanceIDStr); err != nil {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("invalid instance id: %v", err))
		}

		instanceIDInterface := primaryKeyValue.Interface()

		allowed, err := m.App.Panel.PermissionChecker.HasInstanceDeletePermission(m.App.Name, m.Name, instanceIDInterface, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("you are not allowed to delete this instance"))
		}

		err = m.App.Panel.ORM.DeleteInstance(m.PTR, instanceIDInterface)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		return http.StatusFound, m.GetLink()
	}
}

func (m *Model) GetInstanceViewHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := m.App.Panel.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		primaryKeyValue := reflect.New(m.PrimaryKeyType).Elem()
		if err := utils.SetStringsAsType(primaryKeyValue, instanceIDStr); err != nil {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("invalid instance id: %v", err))
		}

		instanceIDInterface := primaryKeyValue.Interface()

		allowed, err := m.App.Panel.PermissionChecker.HasInstanceReadPermission(m.App.Name, m.Name, instanceIDInterface, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("you are not allowed to view this instance"))
		}

		apps, err := GetAppsWithReadPermissions(m.App.Panel, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		var fieldsToFetch []string
		for _, fieldConfig := range m.Fields {
			if fieldConfig.IncludeInInstanceView {
				fieldsToFetch = append(fieldsToFetch, fieldConfig.Name)
			}
		}

		instanceData, err := m.App.Panel.ORM.FetchInstanceOnlyFields(m.PTR, instanceIDInterface, fieldsToFetch)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		html, err := m.App.Panel.Config.Renderer.RenderTemplate("instance.html", map[string]interface{}{
			"model":    m,
			"apps":     apps,
			"instance": instanceData,
		})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		return http.StatusOK, html
	}
}

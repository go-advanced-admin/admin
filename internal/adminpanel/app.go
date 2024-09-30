package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/utils"
	"net/http"
	"reflect"
)

type App struct {
	Name        string
	DisplayName string
	Models      map[string]*Model
	Panel       *AdminPanel
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
		name = modelType.Name()
	}

	if !utils.IsURLSafe(name) {
		return nil, fmt.Errorf("admin model '%s' name is not URL safe", name)
	}

	var displayName string
	displayNamer, ok := model.(AdminModelDisplayNameInterface)
	if ok {
		displayName = displayNamer.AdminDisplayName()
	} else {
		displayName = utils.HumanizeName(name)
	}

	if _, exists := a.Models[name]; exists {
		return nil, fmt.Errorf("admin model '%s' already exists in app '%s'. Models cannot be registered more than once", name, a.Name)
	}
	a.Models[name] = &Model{Name: name, DisplayName: displayName, PTR: model, App: a}
	return a.Models[name], nil
}

func (a *App) GetHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		allowed, err := a.Panel.PermissionChecker.HasAppReadPermission(a.Name, data)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		if !allowed {
			return http.StatusForbidden, "Forbidden"
		}

		models, err := GetModelsWithReadPermissions(a, data)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		html, err := a.Panel.Config.Renderer.RenderTemplate("app.html", map[string]interface{}{"app": a, "models": models})
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusOK, html
	}
}

func (a *App) GetLink() string {
	return fmt.Sprintf("/%s", a.Name)
}

func (a *App) GetFullLink() string {
	return a.Panel.Config.GetLink(a.GetLink())
}

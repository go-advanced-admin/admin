package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal"
	"github.com/go-advanced-admin/admin/internal/utils"
	"net/http"
)

type AdminPanel struct {
	Apps              map[string]*App
	AppsSlice         []*App
	PermissionChecker PermissionFunc
	ORM               ORMIntegrator
	Web               WebIntegrator
	Config            AdminConfig
}

func NewAdminPanel(orm ORMIntegrator, web WebIntegrator, permissionsCheck PermissionFunc, config *AdminConfig) *AdminPanel {
	if config == nil {
		config = &DefaultAdminConfig
	}
	admin := AdminPanel{
		Apps:              make(map[string]*App),
		AppsSlice:         make([]*App, 0),
		PermissionChecker: permissionsCheck,
		ORM:               orm,
		Web:               web,
		Config:            *config,
	}

	admin.Config.Renderer.RegisterDefaultTemplates(internal.TemplateFiles)
	admin.Config.Renderer.RegisterDefaultAssets(internal.AssetsFiles)
	admin.Config.Renderer.RegisterLinkFunc(admin.Config.GetLink)
	admin.Config.Renderer.RegisterAssetsFunc(admin.Config.GetAssetLink)

	web.ServeAssets(config.AssetsPrefix, config.Renderer)
	web.HandleRoute("GET", config.GetPrefix(), admin.GetHandler())

	return &admin
}

func (ap *AdminPanel) GetHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		allowed, err := ap.PermissionChecker.HasReadPermission(data)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		if !allowed {
			return http.StatusForbidden, "Forbidden"
		}

		apps, err := GetAppsWithReadPermissions(ap, data)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		html, err := ap.Config.Renderer.RenderTemplate("root.html", map[string]interface{}{"admin": ap, "apps": apps})
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusOK, html
	}
}

func (ap *AdminPanel) RegisterApp(name, displayName string) (*App, error) {
	if _, exists := ap.Apps[name]; exists {
		return nil, fmt.Errorf("admin app '%s' already exists. Apps cannot be registered more than once", name)
	}

	if !utils.IsURLSafe(name) {
		return nil, fmt.Errorf("admin app name '%s' is not URL safe", name)
	}

	app := &App{Name: name, DisplayName: displayName, Models: make(map[string]*Model), ModelsSlice: make([]*Model, 0), Panel: ap}
	ap.Apps[name] = app
	ap.AppsSlice = append(ap.AppsSlice, app)
	ap.Web.HandleRoute("GET", ap.Config.GetPrefix()+app.GetLink(), app.GetHandler())
	return ap.Apps[name], nil
}

func (ap *AdminPanel) GetFullLink() string {
	return ap.Config.GetLink("")
}

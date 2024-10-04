package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal"
	"github.com/go-advanced-admin/admin/internal/logging"
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

func (ap *AdminPanel) GetLogEntries(ctx interface{}, maxCount uint) []*logging.LogEntry {
	if ap.Config.LogStore == nil {
		return []*logging.LogEntry{}
	}
	entries, err := ap.Config.LogStore.GetLogEntries()
	if err != nil {
		return []*logging.LogEntry{}
	}
	entries = entries[:min(uint(len(entries)), maxCount)]
	permissibleEntries := make([]*logging.LogEntry, 0)
	for _, entry := range entries {
		allowed, err := ap.PermissionChecker.HasLogViewPermission(ctx, entry.ID)
		if err != nil {
			continue
		}
		if allowed {
			permissibleEntries = append(permissibleEntries, entry)
		}
	}
	return permissibleEntries
}

func (ap *AdminPanel) CreateViewLog(ctx interface{}) error {
	return ap.Config.CreateLog(ctx, logging.LogStoreLevelPanelView, "", nil, "", "")
}

func (ap *AdminPanel) GetORM() ORMIntegrator {
	return ap.ORM
}

func NewAdminPanel(orm ORMIntegrator, web WebIntegrator, permissionsCheck PermissionFunc, config *AdminConfig) (*AdminPanel, error) {
	if orm == nil {
		return nil, fmt.Errorf("orm integrator cannot be nil")
	}
	if web == nil {
		return nil, fmt.Errorf("web integrator cannot be nil")
	}
	if permissionsCheck == nil {
		return nil, fmt.Errorf("permissions check function cannot be nil")
	}
	if config == nil {
		config = NewDefaultAdminConfig()
	}
	admin := AdminPanel{
		Apps:              make(map[string]*App),
		AppsSlice:         make([]*App, 0),
		PermissionChecker: permissionsCheck,
		ORM:               orm,
		Web:               web,
		Config:            *config,
	}

	admin.Config.Renderer.RegisterDefaultTemplates(internal.TemplateFiles, "templates/")
	admin.Config.Renderer.RegisterDefaultAssets(internal.AssetsFiles, "assets/")
	admin.Config.Renderer.RegisterLinkFunc(admin.Config.GetLink)
	admin.Config.Renderer.RegisterAssetsFunc(admin.Config.GetAssetLink)

	components := []string{"page.html"}
	pages := []string{"root", "app", "model", "instance", "edit_instance", "new_instance"}

	for _, page := range pages {
		err := admin.Config.Renderer.RegisterCompositeDefaultTemplate(page, append([]string{page + ".html"}, components...)...)
		if err != nil {
			return nil, err
		}
	}

	web.ServeAssets(config.AssetsPrefix, config.Renderer)
	web.HandleRoute("GET", config.GetPrefix(), admin.GetHandler())

	return &admin, nil
}

func (ap *AdminPanel) GetHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		allowed, err := ap.PermissionChecker.HasReadPermission(data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("forbidden"))
		}

		apps, err := GetAppsWithReadPermissions(ap, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		html, err := ap.Config.Renderer.RenderTemplate("root", map[string]interface{}{
			"admin":       ap,
			"apps":        apps,
			"navBarItems": ap.Config.GetNavBarItems(data),
			"logs":        ap.GetLogEntries(data, 20),
		})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		err = ap.CreateViewLog(data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		return http.StatusOK, html
	}
}

func (ap *AdminPanel) RegisterApp(name, displayName string, orm ORMIntegrator) (*App, error) {
	if _, exists := ap.Apps[name]; exists {
		return nil, fmt.Errorf("admin app '%s' already exists. Apps cannot be registered more than once", name)
	}

	if !utils.IsURLSafe(name) {
		return nil, fmt.Errorf("admin app name '%s' is not URL safe", name)
	}

	app := &App{Name: name, DisplayName: displayName, Models: make(map[string]*Model), ModelsSlice: make([]*Model, 0), Panel: ap, ORM: orm}
	ap.Apps[name] = app
	ap.AppsSlice = append(ap.AppsSlice, app)
	ap.Web.HandleRoute("GET", ap.Config.GetPrefix()+app.GetLink(), app.GetHandler())
	return ap.Apps[name], nil
}

func (ap *AdminPanel) GetFullLink() string {
	return ap.Config.GetLink("")
}

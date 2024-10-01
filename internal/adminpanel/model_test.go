package adminpanel

import (
	"embed"
	"net/http"
	"testing"
)

type TestModel struct {
	ID   uint
	Name string
}

type TestWebIntegrator struct{}

func (m *TestWebIntegrator) HandleRoute(method, path string, handler HandlerFunc) {}
func (m *TestWebIntegrator) ServeAssets(prefix string, renderer TemplateRenderer) {}
func (m *TestWebIntegrator) GetQueryParam(ctx interface{}, name string) string {
	queryParams := ctx.(map[string]string)
	return queryParams[name]
}

type TestORMIntegrator struct{}

func (to *TestORMIntegrator) FetchInstances(model interface{}) (interface{}, error) {
	return nil, nil
}

func (to *TestORMIntegrator) FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error) {
	return []TestModel{{ID: 1, Name: "Test"}}, nil
}

type TestTemplateRenderer struct{}

func (tt *TestTemplateRenderer) RenderTemplate(name string, data map[string]interface{}) (string, error) {
	return "Rendered HTML", nil
}

func (tt *TestTemplateRenderer) RegisterDefaultTemplates(templates embed.FS) {}

func (tt *TestTemplateRenderer) RegisterDefaultData(data map[string]interface{}) {}

func (tt *TestTemplateRenderer) AddCustomTemplate(name string, tmplText string) error { return nil }

func (tt *TestTemplateRenderer) RegisterDefaultAssets(assets embed.FS) {}

func (tt *TestTemplateRenderer) AddCustomAsset(name string, asset []byte) {}

func (tt *TestTemplateRenderer) GetAsset(name string) ([]byte, error) { return nil, nil }

func (tt *TestTemplateRenderer) RegisterLinkFunc(func(string) string) {}

func (tt *TestTemplateRenderer) RegisterAssetsFunc(func(string) string) {}

func TestModel_GetLink(t *testing.T) {
	app := &App{Name: "App", Panel: &AdminPanel{Config: AdminConfig{Prefix: "admin"}}}
	model := Model{Name: "TestModel", App: app}

	expectedLink := "/App/TestModel"
	expectedFullLink := "/admin/App/TestModel"

	if link := model.GetLink(); link != expectedLink {
		t.Errorf("expected %s, got %s", expectedLink, link)
	}

	if fullLink := model.GetFullLink(); fullLink != expectedFullLink {
		t.Errorf("expected %s, got %s", expectedFullLink, fullLink)
	}
}

func TestGetPrimaryKeyGetter(t *testing.T) {
	model := &TestModel{}
	getter, err := GetPrimaryKeyGetter(model)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	id := getter(model)
	if id != uint(0) {
		t.Errorf("expected 0, got %v", id)
	}
}

func TestModel_GetViewHandler(t *testing.T) {
	permissionChecker := PermissionFunc(func(req PermissionRequest, data interface{}) (bool, error) {
		return true, nil
	})

	app := &App{
		Name: "App",
		Panel: &AdminPanel{
			ORM: &TestORMIntegrator{},
			Web: &TestWebIntegrator{},
			Config: AdminConfig{
				DefaultInstancesPerPage: 10,
				Renderer:                &TestTemplateRenderer{},
			},
			PermissionChecker: permissionChecker, // Ensure PermissionChecker is set
		},
	}

	model := Model{
		Name: "TestModel",
		App:  app,
		Fields: []FieldConfig{
			{Name: "Name", IncludeInListFetch: true},
			{Name: "ID", IncludeInListFetch: true},
		},
		PrimaryKeyGetter: func(instance interface{}) interface{} { return instance.(TestModel).ID },
	}

	handler := model.GetViewHandler()
	code, html := handler(map[string]string{"page": "1", "perPage": "10"})

	// Validate response
	if code != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, code)
	}

	if html != "Rendered HTML" {
		t.Errorf("expected rendered HTML, got %v", html)
	}
}

func TestFilterInstancesByPermission(t *testing.T) {
	permissionChecker := PermissionFunc(func(req PermissionRequest, data interface{}) (bool, error) {
		return data == nil, nil
	})

	model := Model{
		PrimaryKeyGetter: func(instance interface{}) interface{} { return instance.(TestModel).ID },
		App: &App{
			Name: "App",
			Panel: &AdminPanel{
				PermissionChecker: permissionChecker,
			},
		},
	}

	instances := []TestModel{
		{ID: 1, Name: "Instance1"},
		{ID: 2, Name: "Instance2"},
	}

	filtered, err := filterInstancesByPermission(instances, &model, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filtered) != 2 {
		t.Errorf("expected 2 instances, got %d", len(filtered))
	}

	filtered, err = filterInstancesByPermission(instances, &model, "deny")
	if len(filtered) != 0 {
		t.Errorf("expected 0 instances, got %d", len(filtered))
	}
}

package adminpanel

import (
	"net/http"
	"testing"
)

type TestModel struct {
	ID   uint
	Name string
}

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
	panel, err := NewAdminPanel(&MockORMIntegrator{}, &MockWebIntegrator{}, MockPermissionFunc, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	testApp, err := panel.RegisterApp("TestApp", "Test App")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model, err := testApp.RegisterModel(&TestModel{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	handler := model.GetViewHandler()
	code, _ := handler(map[string]string{"page": "1", "perPage": "10"})

	if code != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, code)
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

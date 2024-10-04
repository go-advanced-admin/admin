package adminpanel

import (
	"net/http"
	"testing"
)

type TestModel struct {
	ID   uint
	Name string
}

type TestModelWithID struct {
	IDValue uint
	Name    string
}

func (m *TestModelWithID) AdminGetID() interface{} {
	return m.IDValue
}

func TestPrimaryKeyGetter_InterfaceImplemented(t *testing.T) {
	model := &TestModelWithID{IDValue: 42}
	getter, err := GetPrimaryKeyGetter(model)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	id := getter(model)
	if id != uint(42) {
		t.Errorf("expected ID 42, got %v", id)
	}
}

func TestGetPrimaryKeyGetter_IDField(t *testing.T) {
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
	panel, err := NewMockAdminPanel()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testApp, err := panel.RegisterApp("TestApp", "Test App", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model, err := testApp.RegisterModel(&TestModel{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name         string
		query        map[string]string
		expectedCode int
	}{
		{"Valid pagination", map[string]string{"page": "1", "perPage": "10"}, http.StatusOK},
		{"Invalid pagination strings", map[string]string{"page": "abc", "perPage": "-1"}, http.StatusOK},
		{"Out of range pagination", map[string]string{"page": "10000", "perPage": "10000"}, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := model.GetViewHandler()
			code, _ := handler(tt.query)
			if code != uint(tt.expectedCode) {
				t.Errorf("expected %v, got %v", tt.expectedCode, code)
			}
		})
	}
}

func TestFilterInstancesByPermission_SliceCheck(t *testing.T) {
	model := Model{
		PrimaryKeyGetter: func(instance interface{}) interface{} {
			if instance, ok := instance.(TestModel); ok {
				return instance.ID
			}
			return nil
		},
		App: &App{
			Name: "App",
			Panel: &AdminPanel{
				PermissionChecker: MockPermissionFunc,
			},
		},
	}

	_, err := filterInstancesByPermission("not a slice", &model, nil)
	if err == nil {
		t.Errorf("expected an error when passing non-slice type")
	}

	instances := []TestModel{
		{ID: 1, Name: "Instance1"},
		{ID: 2, Name: "Instance2"},
	}
	filtered, err := filterInstancesByPermission(instances, &model, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(filtered) != len(instances) {
		t.Errorf("expected %d instances, got %d", len(instances), len(filtered))
	}
}

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

	expectedLink := "/a/App/TestModel"
	expectedFullLink := "/admin/a/App/TestModel"

	if link := model.GetLink(); link != expectedLink {
		t.Errorf("expected %s, got %s", expectedLink, link)
	}

	if fullLink := model.GetFullLink(); fullLink != expectedFullLink {
		t.Errorf("expected %s, got %s", expectedFullLink, fullLink)
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

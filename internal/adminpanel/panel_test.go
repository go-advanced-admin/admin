package adminpanel

import (
	"testing"
)

func TestAdminPanel_RegisterApp(t *testing.T) {
	panel, err := NewAdminPanel(&MockORMIntegrator{}, &MockWebIntegrator{}, MockPermissionFunc, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	app, err := panel.RegisterApp("TestApp", "Test App")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if app == nil {
		t.Fatalf("expected app to be registered, got nil")
	}

	_, err = panel.RegisterApp("TestApp", "Test App Duplicate")
	if err == nil {
		t.Fatalf("expected an error when registering the same app twice")
	}
}

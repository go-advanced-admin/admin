package adminpanel

import (
	"errors"
	"net/http"
	"testing"
)

func PanelMockPermissionFunc(req PermissionRequest, ctx interface{}) (bool, error) {
	if req.Action != nil && *req.Action == ReadAction {
		return true, nil
	}
	return false, nil
}

func TestNewAdminPanel(t *testing.T) {
	_, err := NewAdminPanel(nil, &MockWebIntegrator{}, PanelMockPermissionFunc, nil)
	if err == nil || err.Error() != "orm integrator cannot be nil" {
		t.Fatalf("expected error for nil ORM, got %v", err)
	}

	_, err = NewAdminPanel(&MockORMIntegrator{}, nil, PanelMockPermissionFunc, nil)
	if err == nil || err.Error() != "web integrator cannot be nil" {
		t.Fatalf("expected error for nil Web, got %v", err)
	}

	_, err = NewAdminPanel(&MockORMIntegrator{}, &MockWebIntegrator{}, nil, nil)
	if err == nil || err.Error() != "permissions check function cannot be nil" {
		t.Fatalf("expected error for nil PermissionFunc, got %v", err)
	}

	_, err = NewMockAdminPanel()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAdminPanel_GetHandler(t *testing.T) {
	panel, err := NewMockAdminPanel()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	handler := panel.GetHandler()
	status, _ := handler(nil)

	if status != http.StatusOK {
		t.Errorf("expected status OK, got %d", status)
	}

	deniedPermissionFunc := PermissionFunc(func(req PermissionRequest, ctx interface{}) (bool, error) {
		return false, nil
	})
	panel.PermissionChecker = deniedPermissionFunc

	status, _ = handler(nil)
	if status != http.StatusForbidden {
		t.Errorf("expected status Forbidden, got %d", status)
	}

	// Test handler with permission check error
	errorPermissionFunc := PermissionFunc(func(req PermissionRequest, ctx interface{}) (bool, error) {
		return false, errors.New("permission check error")
	})
	panel.PermissionChecker = errorPermissionFunc

	status, body := handler(nil)
	if status != http.StatusInternalServerError || body != "permission check error" {
		t.Errorf("expected status InternalServerError and error message, got %d and %s", status, body)
	}
}

func TestAdminPanel_RegisterApp(t *testing.T) {
	panel, err := NewMockAdminPanel()
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
	if err == nil || err.Error() != "admin app 'TestApp' already exists. Apps cannot be registered more than once" {
		t.Fatalf("expected error when registering the same app twice, got %v", err)
	}

	_, err = panel.RegisterApp("Unsafe App!", "Unsafe App")
	if err == nil || err.Error() != "admin app name 'Unsafe App!' is not URL safe" {
		t.Fatalf("expected error for unsafe app name, got %v", err)
	}
}

func TestAdminPanel_GetFullLink(t *testing.T) {
	panel, err := NewMockAdminPanel()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := panel.Config.GetLink("")
	if link := panel.GetFullLink(); link != expected {
		t.Errorf("expected %s, got %s", expected, link)
	}
}

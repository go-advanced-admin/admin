package adminpanel

import (
	"testing"
)

func TestPermissionFunc_HasPermission(t *testing.T) {
	permFunc := PermissionFunc(MockPermissionFunc)

	action := ReadAction
	req := PermissionRequest{Action: &action}

	allowed, err := permFunc.HasPermission(req, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !allowed {
		t.Fatalf("expected permission to be allowed")
	}
}

func TestPermissionFunc_OtherActions(t *testing.T) {
	permFunc := PermissionFunc(MockPermissionFunc)

	appName := "App"
	modelName := "Model"

	var (
		varReadAction   = ReadAction
		varCreateAction = CreateAction
		varUpdateAction = UpdateAction
		varDeleteAction = DeleteAction
	)

	tests := []struct {
		name          string
		action        Action
		permissionReq PermissionRequest
		expectError   bool
		expectAllowed bool
	}{
		{
			name:          "Read Permission",
			action:        ReadAction,
			permissionReq: PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &varReadAction},
			expectError:   false,
			expectAllowed: true,
		},
		{
			name:          "Create Permission",
			action:        CreateAction,
			permissionReq: PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &varCreateAction},
			expectError:   false,
			expectAllowed: true,
		},
		{
			name:          "Update Permission",
			action:        UpdateAction,
			permissionReq: PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &varUpdateAction},
			expectError:   false,
			expectAllowed: true,
		},
		{
			name:          "Delete Permission",
			action:        DeleteAction,
			permissionReq: PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &varDeleteAction},
			expectError:   false,
			expectAllowed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			allowed, err := permFunc.HasPermission(test.permissionReq, nil)
			if (err != nil) != test.expectError {
				t.Fatalf("expected error %v, got %v", test.expectError, err)
			}
			if allowed != test.expectAllowed {
				t.Fatalf("expected %v, got %v", test.expectAllowed, allowed)
			}
		})
	}
}

func TestPermissionFuncs(t *testing.T) {
	permFunc := PermissionFunc(MockPermissionFunc)

	tests := []struct {
		name          string
		permissionReq func(appName, modelName string, data interface{}) (bool, error)
		appName       string
		modelName     string
		expectAllowed bool
	}{
		{
			name:          "HasReadPermission",
			permissionReq: func(_, _ string, data interface{}) (bool, error) { return permFunc.HasReadPermission(data) },
			appName:       "",
			modelName:     "",
			expectAllowed: true,
		},
		{
			name: "HasAppReadPermission",
			permissionReq: func(appName, _ string, data interface{}) (bool, error) {
				return permFunc.HasAppReadPermission(appName, data)
			},
			appName:       "App",
			modelName:     "",
			expectAllowed: true,
		},
		{
			name:          "HasModelReadPermission",
			permissionReq: permFunc.HasModelReadPermission,
			appName:       "App",
			modelName:     "Model",
			expectAllowed: true,
		},
		{
			name:          "HasModelCreatePermission",
			permissionReq: permFunc.HasModelCreatePermission,
			appName:       "App",
			modelName:     "Model",
			expectAllowed: true,
		},
		{
			name:          "HasModelUpdatePermission",
			permissionReq: permFunc.HasModelUpdatePermission,
			appName:       "App",
			modelName:     "Model",
			expectAllowed: true,
		},
		{
			name:          "HasModelDeletePermission",
			permissionReq: permFunc.HasModelDeletePermission,
			appName:       "App",
			modelName:     "Model",
			expectAllowed: true,
		},
		{
			name: "HasInstanceReadPermission",
			permissionReq: func(appName, modelName string, data interface{}) (bool, error) {
				return permFunc.HasInstanceReadPermission(appName, modelName, nil, data)
			},
			appName:       "App",
			modelName:     "Model",
			expectAllowed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			allowed, err := test.permissionReq(test.appName, test.modelName, nil)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if allowed != test.expectAllowed {
				t.Fatalf("expected %v, got %v", test.expectAllowed, allowed)
			}
		})
	}
}

func TestPermissionsRetrieval(t *testing.T) {
	adminPanel, err := NewMockAdminPanel()
	if err != nil {
		t.Fatalf("failed to create mock admin panel: %v", err)
	}

	app, err := adminPanel.RegisterApp("TestApp", "Test App", nil)
	if err != nil {
		t.Fatalf("failed to register app: %v", err)
	}

	model := &TestModel1{}
	_, err = app.RegisterModel(model, nil)
	if err != nil {
		t.Fatalf("failed to register model: %v", err)
	}

	models, err := GetModelsWithReadPermissions(app, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}

	modelPerm := models[0]["permissions"].(Permissions)
	if !modelPerm.Read || !modelPerm.Create || !modelPerm.Update || !modelPerm.Delete {
		t.Fatalf("expected all permissions to be true, got %v", modelPerm)
	}

	apps, err := GetAppsWithReadPermissions(adminPanel, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}
}

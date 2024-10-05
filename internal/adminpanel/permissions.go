package adminpanel

// Action represents an action type for permissions.
type Action string

const (
	// ReadAction represents read permissions.
	ReadAction Action = "read"
	// CreateAction represents create permissions.
	CreateAction Action = "create"
	// UpdateAction represents update permissions.
	UpdateAction Action = "update"
	// DeleteAction represents delete permissions.
	DeleteAction Action = "delete"
	// LogViewAction represents log viewing permissions.
	LogViewAction Action = "log_view"
)

// PermissionRequest represents a request to check permissions for a specific action.
type PermissionRequest struct {
	AppName    *string
	ModelName  *string
	InstanceID interface{}
	Action     *Action
}

// Permissions holds the permissions for a specific operation.
type Permissions struct {
	Read   bool
	Create bool
	Update bool
	Delete bool
}

// PermissionFunc defines a function type for checking permissions.
type PermissionFunc func(PermissionRequest, interface{}) (bool, error)

// HasLogViewPermission checks if the user has permission to view logs.
func (p PermissionFunc) HasLogViewPermission(data interface{}, logID interface{}) (bool, error) {
	action := LogViewAction
	return p(PermissionRequest{Action: &action, InstanceID: logID}, data)
}

// HasPermission checks if the user has the specified permission.
func (p PermissionFunc) HasPermission(r PermissionRequest, data interface{}) (bool, error) {
	return p(r, data)
}

// HasReadPermission checks if the user has read permission for the admin panel.
func (p PermissionFunc) HasReadPermission(data interface{}) (bool, error) {
	action := ReadAction
	return p(PermissionRequest{Action: &action}, data)
}

// HasAppReadPermission checks if the user has read permission for the specified app.
func (p PermissionFunc) HasAppReadPermission(appName string, data interface{}) (bool, error) {
	action := ReadAction
	permissionRequest := PermissionRequest{AppName: &appName, Action: &action}
	return p(permissionRequest, data)
}

// HasModelReadPermission checks if the user has read permission for the specified model.
func (p PermissionFunc) HasModelReadPermission(appName string, modelName string, data interface{}) (bool, error) {
	action := ReadAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action}
	return p(permissionRequest, data)
}

// HasModelCreatePermission checks if the user has create permission for the specified model.
func (p PermissionFunc) HasModelCreatePermission(appName string, modelName string, data interface{}) (bool, error) {
	action := CreateAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action}
	return p(permissionRequest, data)
}

// HasModelUpdatePermission checks if the user has update permission for the specified model.
func (p PermissionFunc) HasModelUpdatePermission(appName string, modelName string, data interface{}) (bool, error) {
	action := UpdateAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action}
	return p(permissionRequest, data)
}

// HasModelDeletePermission checks if the user has delete permission for the specified model.
func (p PermissionFunc) HasModelDeletePermission(appName string, modelName string, data interface{}) (bool, error) {
	action := DeleteAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action}
	return p(permissionRequest, data)
}

// HasInstanceReadPermission checks if the user has read permission for the specified instance.
func (p PermissionFunc) HasInstanceReadPermission(appName, modelName string, instanceID interface{}, data interface{}) (bool, error) {
	action := ReadAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action, InstanceID: instanceID}
	return p(permissionRequest, data)
}

// HasInstanceUpdatePermission checks if the user has update permission for the specified instance.
func (p PermissionFunc) HasInstanceUpdatePermission(appName, modelName string, instanceID interface{}, data interface{}) (bool, error) {
	action := UpdateAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action, InstanceID: instanceID}
	return p(permissionRequest, data)
}

// HasInstanceDeletePermission checks if the user has delete permission for the specified instance.
func (p PermissionFunc) HasInstanceDeletePermission(appName, modelName string, instanceID interface{}, data interface{}) (bool, error) {
	action := DeleteAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action, InstanceID: instanceID}
	return p(permissionRequest, data)
}

// GetModelsWithReadPermissions returns models for which the user has read permissions.
func GetModelsWithReadPermissions(app *App, data interface{}) ([]map[string]interface{}, error) {
	modelsSlice := make([]map[string]interface{}, 0)

	for _, model := range app.ModelsSlice {
		modelMap := make(map[string]interface{})
		modelReadAllowed, err := app.Panel.PermissionChecker.HasModelReadPermission(app.Name, model.Name, data)
		if err != nil {
			return nil, err
		}
		if !modelReadAllowed {
			continue
		}
		modelMap["model"] = model

		createAllowed, err := app.Panel.PermissionChecker.HasModelCreatePermission(app.Name, model.Name, data)
		if err != nil {
			return nil, err
		}

		updateAllowed, err := app.Panel.PermissionChecker.HasModelUpdatePermission(app.Name, model.Name, data)
		if err != nil {
			return nil, err
		}

		deleteAllowed, err := app.Panel.PermissionChecker.HasModelDeletePermission(app.Name, model.Name, data)
		if err != nil {
			return nil, err
		}

		permissions := Permissions{
			Read:   modelReadAllowed,
			Create: createAllowed,
			Update: updateAllowed,
			Delete: deleteAllowed,
		}
		modelMap["permissions"] = permissions
		modelsSlice = append(modelsSlice, modelMap)
	}

	return modelsSlice, nil
}

// GetAppsWithReadPermissions returns apps for which the user has read permissions.
func GetAppsWithReadPermissions(panel *AdminPanel, data interface{}) ([]map[string]interface{}, error) {
	apps := make([]map[string]interface{}, 0)
	for _, app := range panel.AppsSlice {
		appMap := make(map[string]interface{})
		readAllowed, err := panel.PermissionChecker.HasAppReadPermission(app.Name, data)
		if err != nil {
			return nil, err
		}
		if !readAllowed {
			continue
		}
		appMap["app"] = app
		modelsSlice, err := GetModelsWithReadPermissions(app, data)
		if err != nil {
			return nil, err
		}

		appMap["models"] = modelsSlice
		apps = append(apps, appMap)
	}
	return apps, nil
}

package adminpanel

type Action string

const (
	ReadAction   Action = "read"
	CreateAction Action = "create"
	UpdateAction Action = "update"
	DeleteAction Action = "delete"
)

type PermissionRequest struct {
	AppName    *string
	ModelName  *string
	InstanceID interface{}
	Action     *Action
}

type Permissions struct {
	Read   bool
	Create bool
	Update bool
	Delete bool
}

type PermissionFunc func(PermissionRequest, interface{}) (bool, error)

func (p PermissionFunc) HasPermission(r PermissionRequest, data interface{}) (bool, error) {
	return p(r, data)
}

func (p PermissionFunc) HasReadPermission(data interface{}) (bool, error) {
	action := ReadAction
	return p(PermissionRequest{Action: &action}, data)
}

func (p PermissionFunc) HasAppReadPermission(appName string, data interface{}) (bool, error) {
	action := ReadAction
	permissionRequest := PermissionRequest{AppName: &appName, Action: &action}
	return p(permissionRequest, data)
}

func (p PermissionFunc) HasModelReadPermission(appName string, modelName string, data interface{}) (bool, error) {
	action := ReadAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action}
	return p(permissionRequest, data)
}

func (p PermissionFunc) HasModelCreatePermission(appName string, modelName string, data interface{}) (bool, error) {
	action := CreateAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action}
	return p(permissionRequest, data)
}

func (p PermissionFunc) HasModelUpdatePermission(appName string, modelName string, data interface{}) (bool, error) {
	action := UpdateAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action}
	return p(permissionRequest, data)
}

func (p PermissionFunc) HasModelDeletePermission(appName string, modelName string, data interface{}) (bool, error) {
	action := DeleteAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action}
	return p(permissionRequest, data)
}

func (p PermissionFunc) HasInstanceReadPermission(appName, modelName string, instanceID interface{}, data interface{}) (bool, error) {
	action := ReadAction
	permissionRequest := PermissionRequest{AppName: &appName, ModelName: &modelName, Action: &action, InstanceID: instanceID}
	return p(permissionRequest, data)
}

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

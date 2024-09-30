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
	InstanceID *interface{}
	Action     *Action
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

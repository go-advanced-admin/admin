package adminpanel

type PermissionRequest struct {
	AppName    *string
	ModelName  *string
	InstanceID *interface{}
	Action     *string
}

type PermissionFunc func(PermissionRequest, interface{}) (bool, error)

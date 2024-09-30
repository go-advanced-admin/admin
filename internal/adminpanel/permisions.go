package adminpanel

type PermissionRequest struct {
	AppName    *string
	ModelName  *string
	InstanceID *interface{}
	Action     *string
}

type PermissionFunc func(PermissionRequest, interface{}) (bool, error)

func (p PermissionFunc) HasPermission(r PermissionRequest, data interface{}) (bool, error) {
	return p(r, data)
}

func (p PermissionFunc) HasViewPermission(data interface{}) (bool, error) {
	action := "view"
	return p(PermissionRequest{Action: &action}, data)
}

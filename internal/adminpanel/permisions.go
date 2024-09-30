package adminpanel

type PermissionRequest struct {
	AppName    *string
	ModelName  *string
	InstanceID *interface{}
	Action     *string
}

type PermissionFunc func(PermissionRequest, interface{}) (bool, error)

func (p PermissionFunc) HasPermission(r PermissionRequest, instance interface{}) (bool, error) {
	return p(r, instance)
}

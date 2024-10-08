package adminpanel

func MockPermissionFunc(_ PermissionRequest, _ interface{}) (bool, error) {
	return true, nil
}

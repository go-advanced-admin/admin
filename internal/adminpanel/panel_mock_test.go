package adminpanel

func NewMockAdminPanel() (*AdminPanel, error) {
	return NewAdminPanel(&MockORMIntegrator{}, &MockWebIntegrator{}, MockPermissionFunc, nil)
}

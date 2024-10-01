package adminpanel

type MockORMIntegrator struct{}

func (m *MockORMIntegrator) FetchInstances(model interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockORMIntegrator) FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error) {
	return []interface{}{}, nil
}

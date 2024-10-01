package adminpanel

type MockORMIntegrator struct{}

func (m *MockORMIntegrator) FetchInstances(interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockORMIntegrator) FetchInstancesOnlyFields(interface{}, []string) (interface{}, error) {
	return []interface{}{}, nil
}

func (m *MockORMIntegrator) FetchInstancesOnlyFieldWithSearch(interface{}, []string, string, []string) (interface{}, error) {
	return []interface{}{}, nil
}

package adminpanel

type ORMIntegrator interface {
	FetchInstances(model interface{}) (interface{}, error)
	FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error)
}

type MockORMIntegrator struct{}

func (m *MockORMIntegrator) FetchInstances(model interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockORMIntegrator) FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error) {
	return []interface{}{}, nil
}

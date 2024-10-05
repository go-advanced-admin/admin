package adminpanel

import "reflect"

type MockORMIntegrator struct{}

func (m *MockORMIntegrator) GetPrimaryKeyValue(interface{}) (interface{}, error) { return nil, nil }
func (m *MockORMIntegrator) GetPrimaryKeyType(interface{}) (reflect.Type, error) { return nil, nil }

func (m *MockORMIntegrator) FetchInstances(interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockORMIntegrator) FetchInstancesOnlyFields(interface{}, []string) (interface{}, error) {
	return []interface{}{}, nil
}

func (m *MockORMIntegrator) FetchInstancesOnlyFieldWithSearch(interface{}, []string, string, []string) (interface{}, error) {
	return []interface{}{}, nil
}

func (m *MockORMIntegrator) DeleteInstance(interface{}, interface{}) error {
	return nil
}

func (m *MockORMIntegrator) FetchInstanceOnlyFields(interface{}, interface{}, []string) (interface{}, error) {
	return nil, nil
}

func (m *MockORMIntegrator) FetchInstance(interface{}, interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockORMIntegrator) CreateInstance(interface{}) error {
	return nil
}
func (m *MockORMIntegrator) UpdateInstance(interface{}, interface{}) error {
	return nil
}
func (m *MockORMIntegrator) CreateInstanceOnlyFields(interface{}, []string) error {
	return nil
}
func (m *MockORMIntegrator) UpdateInstanceOnlyFields(interface{}, []string, interface{}) error {
	return nil
}

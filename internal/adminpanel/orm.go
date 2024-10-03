package adminpanel

type ORMIntegrator interface {
	FetchInstances(model interface{}) (interface{}, error)
	FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error)
	FetchInstancesOnlyFieldWithSearch(model interface{}, fields []string, query string, searchFields []string) (interface{}, error)
	DeleteInstance(model interface{}, id interface{}) error
	FetchInstanceOnlyFields(model interface{}, id interface{}, fields []string) (interface{}, error)
	CreateInstance(instance interface{}) error
	UpdateInstance(instance interface{}) error
	CreateInstanceOnlyFields(instance interface{}, fields []string) error
	UpdateInstanceOnlyFields(instance interface{}, fields []string) error
}

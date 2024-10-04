package adminpanel

import "reflect"

type ORMIntegrator interface {
	GetPrimaryKeyValue(model interface{}) (interface{}, error)
	GetPrimaryKeyType(model interface{}) (reflect.Type, error)
	FetchInstances(model interface{}) (interface{}, error)
	FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error)
	FetchInstancesOnlyFieldWithSearch(model interface{}, fields []string, query string, searchFields []string) (interface{}, error)
	DeleteInstance(model interface{}, id interface{}) error
	FetchInstanceOnlyFields(model interface{}, id interface{}, fields []string) (interface{}, error)
	FetchInstance(model interface{}, id interface{}) (interface{}, error)
	CreateInstance(instance interface{}) error
	UpdateInstance(instance interface{}, primaryKey interface{}) error
	CreateInstanceOnlyFields(instance interface{}, fields []string) error
	UpdateInstanceOnlyFields(instance interface{}, fields []string, primaryKey interface{}) error
}

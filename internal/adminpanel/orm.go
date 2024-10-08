package adminpanel

import "reflect"

// ORMIntegrator defines the interface for integrating ORMs with the admin panel.
type ORMIntegrator interface {
	// GetPrimaryKeyValue returns the primary key value of the given model instance.
	GetPrimaryKeyValue(model interface{}) (interface{}, error)

	// GetPrimaryKeyType returns the reflect.Type of the primary key for the model.
	GetPrimaryKeyType(model interface{}) (reflect.Type, error)

	// FetchInstances retrieves all instances of the given model.
	FetchInstances(model interface{}) (interface{}, error)

	// FetchInstancesOnlyFields retrieves instances with only the specified fields.
	FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error)

	// FetchInstancesOnlyFieldWithSearch retrieves instances matching the search query on specified fields.
	FetchInstancesOnlyFieldWithSearch(model interface{}, fields []string, query string, searchFields []string) (interface{}, error)

	// DeleteInstance deletes an instance of the model by its primary key.
	DeleteInstance(model interface{}, id interface{}) error

	// FetchInstanceOnlyFields retrieves a single instance with only the specified fields.
	FetchInstanceOnlyFields(model interface{}, id interface{}, fields []string) (interface{}, error)

	// FetchInstance retrieves a single instance of the model by its primary key.
	FetchInstance(model interface{}, id interface{}) (interface{}, error)

	// CreateInstance creates a new instance of the model.
	CreateInstance(instance interface{}) error

	// UpdateInstance updates an existing instance of the model.
	UpdateInstance(instance interface{}, primaryKey interface{}) error

	// CreateInstanceOnlyFields creates a new instance with only the specified fields.
	CreateInstanceOnlyFields(instance interface{}, fields []string) error

	// UpdateInstanceOnlyFields updates an existing instance with only the specified fields.
	UpdateInstanceOnlyFields(instance interface{}, fields []string, primaryKey interface{}) error
}

package adminpanel

type ORMIntegrator interface {
	FetchInstances(model interface{}) (interface{}, error)
	FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error)
}

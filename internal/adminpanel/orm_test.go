package adminpanel

import (
	"reflect"
	"testing"
)

func TestMockORMIntegrator_FetchInstances(t *testing.T) {
	mockORM := &MockORMIntegrator{}

	model := &TestModel{}
	result, err := mockORM.FetchInstances(model)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("expected result to be nil, got %v", result)
	}
}

func TestMockORMIntegrator_FetchInstancesOnlyFields(t *testing.T) {
	mockORM := &MockORMIntegrator{}

	model := &TestModel{}
	fields := []string{"ID", "Name"}
	result, err := mockORM.FetchInstancesOnlyFields(model, fields)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []interface{}{}) {
		t.Errorf("expected result to be an empty slice, got %v", result)
	}
}

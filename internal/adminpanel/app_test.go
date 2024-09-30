package adminpanel

import (
	"testing"
)

type TestModel1 struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func TestRegisterModel(t *testing.T) {
	panel, err := NewAdminPanel(&MockORMIntegrator{}, &MockWebIntegrator{}, MockPermissionFunc, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testApp, err := panel.RegisterApp("TestApp", "Test App")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testModel := &TestModel1{}

	model, err := testApp.RegisterModel(testModel)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if model == nil {
		t.Fatalf("expected model to be registered, got nil")
	}

	_, err = testApp.RegisterModel(testModel)
	if err == nil {
		t.Fatalf("expected an error when registering the same model twice")
	}
}

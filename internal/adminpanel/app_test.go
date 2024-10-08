package adminpanel

import (
	"errors"
	"net/http"
	"testing"
)

type TestModel1 struct {
	ID   uint   `gorm:"primarykey"`
	Name string `admin:"displayName:Custom Name;listFetch:include"`
}

type TestModel2 struct {
	ID     uint   `gorm:"primarykey"`
	Status string `admin:"listDisplay:exclude;unknownTag:test"`
}

type NonPointerModel struct {
	Name string
}

type CustomNameModel struct {
	ID uint
}

func (m *CustomNameModel) AdminName() string {
	return "CustomModelName"
}

func (m *CustomNameModel) AdminDisplayName() string {
	return "Custom Display Name"
}

type unsafeModelName struct {
	ID uint
}

func (m *unsafeModelName) AdminName() string {
	return "unsafe Model Name"
}

func TestRegisterModel(t *testing.T) {
	createTestApp := func() *App {
		panel, err := NewMockAdminPanel()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testApp, err := panel.RegisterApp("TestApp", "Test App", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return testApp
	}

	t.Run("ValidModel", func(t *testing.T) {
		testApp := createTestApp()
		testModel := &struct {
			ID uint `gorm:"primarykey"`
		}{}
		model, err := testApp.RegisterModel(testModel, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if model == nil {
			t.Fatalf("expected model to be registered, got nil")
		}
		if len(testApp.Models) != 1 {
			t.Fatalf("expected one model in Models map, got %d", len(testApp.Models))
		}
		if len(testApp.ModelsSlice) != 1 {
			t.Fatalf("expected one model in ModelsSlice, got %d", len(testApp.ModelsSlice))
		}
	})

	t.Run("DuplicateModel", func(t *testing.T) {
		testApp := createTestApp()
		testModel := &TestModel1{}
		_, err := testApp.RegisterModel(testModel, nil)
		if err != nil {
			t.Fatalf("expected no error on first registration, got %v", err)
		}
		_, err = testApp.RegisterModel(testModel, nil)
		if err == nil {
			t.Error("expected an error when registering the same model twice")
		}
	})

	t.Run("NonPointerModel", func(t *testing.T) {
		testApp := createTestApp()
		_, err := testApp.RegisterModel(NonPointerModel{}, nil)
		if err == nil {
			t.Error("expected an error when registering a non-pointer model")
		}
	})

	t.Run("PointerToNonStruct", func(t *testing.T) {
		testApp := createTestApp()
		var testString *string
		_, err := testApp.RegisterModel(testString, nil)
		if err == nil {
			t.Error("expected an error when registering a pointer to a non-struct type")
		}
	})

	t.Run("FieldConfiguration", func(t *testing.T) {
		testApp := createTestApp()
		testModel := &TestModel1{}
		model, err := testApp.RegisterModel(testModel, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(model.Fields) < 2 {
			t.Fatalf("expected model to have at least 2 fields, got %d", len(model.Fields))
		}
		if model.Fields[1].DisplayName != "Custom Name" {
			t.Fatalf("expected display name 'Custom Name', got '%s'", model.Fields[1].DisplayName)
		}
		if !model.Fields[1].IncludeInListFetch {
			t.Error("expected field 'Name' to be included in fetch")
		}
	})

	t.Run("UnknownTagKey", func(t *testing.T) {
		testApp := createTestApp()
		_, err := testApp.RegisterModel(&TestModel2{}, nil)
		if err != nil {
			t.Error("expected no error due to unknown tag key")
		}
	})

	t.Run("InvalidTagValueForListDisplay", func(t *testing.T) {
		testApp := createTestApp()
		type InvalidTagValueModel struct {
			ID   uint   `gorm:"primarykey"`
			Name string `admin:"listDisplay:invalid"`
		}

		_, err := testApp.RegisterModel(&InvalidTagValueModel{}, nil)
		if err == nil {
			t.Error("expected an error due to invalid listDisplay tag value")
		}
	})

	t.Run("InvalidTagValueForListFetch", func(t *testing.T) {
		testApp := createTestApp()
		type InvalidTagFetchModel struct {
			ID   uint   `gorm:"primarykey"`
			Name string `admin:"listFetch:invalid"`
		}

		_, err := testApp.RegisterModel(&InvalidTagFetchModel{}, nil)
		if err == nil {
			t.Error("expected an error due to invalid listFetch tag value")
		}
	})

	t.Run("ExplicitlyExcludeFromFetch", func(t *testing.T) {
		testApp := createTestApp()
		type ModelExcludeFetch struct {
			ID   uint   `gorm:"primarykey"`
			Name string `admin:"listFetch:exclude"`
		}

		model, err := testApp.RegisterModel(&ModelExcludeFetch{}, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if model.Fields[1].IncludeInListFetch {
			t.Error("expected field 'Name' to be excluded from fetch")
		}
	})

	t.Run("IncludeInFetchDefaultBehavior", func(t *testing.T) {
		testApp := createTestApp()
		type ModelWithID struct {
			ID     uint
			Status string `admin:"listDisplay:exclude"`
		}

		type ModelWithoutID struct {
			Status string `admin:"listDisplay:exclude"`
		}

		modelWithID, err := testApp.RegisterModel(&ModelWithID{}, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !modelWithID.Fields[0].IncludeInListFetch {
			t.Error("expected ID field to be included in fetch by default")
		}
	})

	t.Run("URLSafetyCheck", func(t *testing.T) {
		testApp := createTestApp()

		_, err := testApp.RegisterModel(&unsafeModelName{}, nil)
		if err == nil {
			t.Error("expected an error due to invalid URL safety of the model's name")
		}
	})

	t.Run("ValidModelWithNameCheck", func(t *testing.T) {
		testApp := createTestApp()

		type safeNameModel struct {
			ID uint
		}

		model, err := testApp.RegisterModel(&safeNameModel{}, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if model.Name != "safeNameModel" {
			t.Errorf("Expected model name to be safe, got: %s", model.Name)
		}
	})

	t.Run("CustomNameAndDisplayName", func(t *testing.T) {
		testApp := createTestApp()
		model, err := testApp.RegisterModel(&CustomNameModel{}, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if model.Name != "CustomModelName" {
			t.Fatalf("expected model name 'CustomModelName', got '%s'", model.Name)
		}
		if model.DisplayName != "Custom Display Name" {
			t.Fatalf("expected display name 'Custom Display Name', got '%s'", model.DisplayName)
		}
	})

	t.Run("PermissionError", func(t *testing.T) {
		panel, err := NewMockAdminPanel()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testApp, err := panel.RegisterApp("TestApp", "Test App", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		panel.PermissionChecker = func(req PermissionRequest, ctx interface{}) (bool, error) {
			return false, errors.New("mock error")
		}
		handlerFunc := testApp.GetHandler()
		status, _ := handlerFunc(nil)
		if status != http.StatusInternalServerError {
			t.Fatalf("expected status 500 for mock error, got '%v'", status)
		}
	})

	t.Run("ForbiddenAccess", func(t *testing.T) {
		panel, err := NewMockAdminPanel()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testApp, err := panel.RegisterApp("TestApp", "Test App", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		panel.PermissionChecker = func(req PermissionRequest, ctx interface{}) (bool, error) {
			return false, nil
		}
		handlerFunc := testApp.GetHandler()
		status, _ := handlerFunc(nil)
		if status != http.StatusForbidden {
			t.Fatalf("expected status 403 for forbidden access, got '%v'", status)
		}
	})

	t.Run("SuccessfulRender", func(t *testing.T) {
		panel, err := NewMockAdminPanel()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testApp, err := panel.RegisterApp("TestApp", "Test App", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		panel.PermissionChecker = MockPermissionFunc
		handlerFunc := testApp.GetHandler()
		status, html := handlerFunc(nil)
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got '%v'", status)
		}
		if html == "" {
			t.Fatal("expected non-empty html content for successful render")
		}
	})
}

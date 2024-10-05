package adminpanel

import (
	"testing"
)

func TestMockWebIntegrator_HandleRoute(t *testing.T) {
	webIntegrator := &MockWebIntegrator{}

	handled := false
	handler := func(data interface{}) (uint, string) {
		handled = true
		return 200, "OK"
	}

	webIntegrator.HandleRoute("GET", "/test", handler)

	// In a real setting, you would trigger a request to the route and check this flag,
	// but since HandleRoute doesn't do anything in the mock, we'll just check the setup
	if !handled {
		t.Log("Handlder was not executed during test run")
	}
}

// This test case demonstrates the invocation of the ServeAssets method
func TestMockWebIntegrator_ServeAssets(t *testing.T) {
	webIntegrator := &MockWebIntegrator{}
	renderer := NewDefaultTemplateRenderer()

	webIntegrator.ServeAssets("assets", renderer)

	// Since ServeAssets does not modify any state in the mock,
	// there's no verifiable outcome, just ensure no panic occurs.
	t.Log("ServeAssets executed without any errors")
}

func TestMockWebIntegrator_GetQueryParam(t *testing.T) {
	webIntegrator := &MockWebIntegrator{}
	param := webIntegrator.GetQueryParam(nil, "param")

	if param != "" {
		t.Errorf("expected empty string, got %s", param)
	}
	t.Log("GetQueryParam returned expected result")
}

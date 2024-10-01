package adminpanel

type MockWebIntegrator struct{}

func (m *MockWebIntegrator) HandleRoute(method, path string, handler HandlerFunc) {}
func (m *MockWebIntegrator) ServeAssets(prefix string, renderer TemplateRenderer) {}
func (m *MockWebIntegrator) GetQueryParam(ctx interface{}, name string) string {
	return ""
}

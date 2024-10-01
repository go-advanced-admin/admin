package adminpanel

type MockWebIntegrator struct{}

func (m *MockWebIntegrator) HandleRoute(string, string, HandlerFunc) {}
func (m *MockWebIntegrator) ServeAssets(string, TemplateRenderer)    {}
func (m *MockWebIntegrator) GetQueryParam(ctx interface{}, name string) string {
	if query, ok := ctx.(map[string]string); ok {
		return query[name]
	}
	return ""
}
func (m *MockWebIntegrator) GetPathParam(ctx interface{}, name string) string {
	if path, ok := ctx.(map[string]string); ok {
		return path[name]
	}
	return ""
}

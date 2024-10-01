package adminpanel

type MockWebIntegrator struct{}

func (m *MockWebIntegrator) HandleRoute(string, string, HandlerFunc) {}
func (m *MockWebIntegrator) ServeAssets(string, TemplateRenderer)    {}
func (m *MockWebIntegrator) GetQueryParam(interface{}, string) string {
	return ""
}

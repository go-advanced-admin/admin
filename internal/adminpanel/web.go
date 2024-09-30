package adminpanel

type HandlerFunc = func(interface{}) (uint, string)

type WebIntegrator interface {
	HandleRoute(method, path string, handler HandlerFunc)
	ServeAssets(prefix string, renderer TemplateRenderer)
	GetQueryParam(ctx interface{}, name string) string
}

type MockWebIntegrator struct{}

func (m *MockWebIntegrator) HandleRoute(method, path string, handler HandlerFunc) {}
func (m *MockWebIntegrator) ServeAssets(prefix string, renderer TemplateRenderer) {}
func (m *MockWebIntegrator) GetQueryParam(ctx interface{}, name string) string {
	return ""
}

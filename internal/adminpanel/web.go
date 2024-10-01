package adminpanel

type HandlerFunc = func(interface{}) (uint, string)

type WebIntegrator interface {
	HandleRoute(method, path string, handler HandlerFunc)
	ServeAssets(prefix string, renderer TemplateRenderer)
	GetQueryParam(ctx interface{}, name string) string
}

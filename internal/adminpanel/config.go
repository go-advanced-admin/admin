package adminpanel

type AdminConfig struct {
	Name         string
	Prefix       string
	Renderer     TemplateRenderer
	AssetsPrefix string
	GroupPrefix  string
}

var DefaultAdminConfig = AdminConfig{
	Name:         "Site Administration",
	Prefix:       "admin",
	AssetsPrefix: "admin-assets",
	Renderer:     NewDefaultTemplateRenderer(),
}

func (c *AdminConfig) GetPrefix() string {
	if c.Prefix == "" {
		return ""
	}
	return "/" + c.Prefix
}

func (c *AdminConfig) GetAssetsPrefix() string {
	if c.AssetsPrefix == "" {
		return ""
	}
	return "/" + c.AssetsPrefix
}

func (c *AdminConfig) GetLink(link string) string {
	return c.GroupPrefix + c.GetPrefix() + link
}

func (c *AdminConfig) GetAssetLink(fileName string) string {
	return c.GroupPrefix + c.GetAssetsPrefix() + "/" + fileName
}

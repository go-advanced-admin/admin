package adminpanel

type AdminConfig struct {
	Name                    string
	Prefix                  string
	Renderer                TemplateRenderer
	AssetsPrefix            string
	GroupPrefix             string
	DefaultInstancesPerPage uint
	NavBarGenerators        []NavBarGenerator
}

var DefaultAdminConfig = NewDefaultAdminConfig()

func NewDefaultAdminConfig() *AdminConfig {
	navBarGens := []NavBarGenerator{
		func(interface{}) NavBarItem { return NavBarItem{Name: "Welcome, User. ", Bold: true} },
		func(interface{}) NavBarItem { return NavBarItem{Name: "View Site", Link: "/"} },
		func(interface{}) NavBarItem { return NavBarItem{Name: "View Site", Link: "/"} },
		func(interface{}) NavBarItem { return NavBarItem{Name: "View Site", Link: "/"} },
	}

	return &AdminConfig{
		Name:                    "Site Administration",
		Prefix:                  "admin",
		AssetsPrefix:            "admin-assets",
		Renderer:                NewDefaultTemplateRenderer(),
		DefaultInstancesPerPage: 10,
		NavBarGenerators:        navBarGens,
	}
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

func (c *AdminConfig) GetNavBarItems(ctx interface{}) []NavBarItem {
	items := make([]NavBarItem, 0)
	for idx, generator := range c.NavBarGenerators {
		item := generator(ctx)

		if idx != len(c.NavBarGenerators)-1 && !item.Bold {
			item.NavBarAppendSlash = true
		}

		html := item.HTML()
		if html != "" {
			items = append(items, item)
		}
	}
	return items
}

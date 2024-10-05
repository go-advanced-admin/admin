package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/logging"
	"github.com/google/uuid"
	"time"
)

// AdminConfig holds configuration settings for the admin panel.
type AdminConfig struct {
	Name                    string
	Prefix                  string
	Renderer                TemplateRenderer
	AssetsPrefix            string
	GroupPrefix             string
	DefaultInstancesPerPage uint
	NavBarGenerators        []NavBarGenerator
	UserFetcher             UserFetchFunction
	LogStore                logging.LogStore
	LogStoreLevel           logging.LogStoreLevel
}

// UserFetchFunction defines a function type for fetching user information from the context.
type UserFetchFunction = func(ctx interface{}) (userID interface{}, repr string, err error)

// DefaultAdminConfig provides default configuration settings for the admin panel.
var DefaultAdminConfig = NewDefaultAdminConfig()

// NewDefaultAdminConfig returns a new AdminConfig with default settings.
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
		LogStore:                logging.NewInMemoryLogStore(100),
		LogStoreLevel:           logging.LogStoreLevelPanelView,
		NavBarGenerators:        navBarGens,
	}
}

// CreateLog creates a log entry using the admin panel's log store.
func (c *AdminConfig) CreateLog(ctx interface{}, action logging.LogStoreLevel, contentType string, objectID interface{}, objectRepr string, message string) error {
	if !c.LogStoreLevel.AssessLevel(action) {
		return nil
	}

	var userId interface{}
	var userRepr string
	var err error
	if c.UserFetcher != nil {
		userId, userRepr, err = c.UserFetcher(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch user: %w", err)
		}
	}
	logEntry := logging.LogEntry{
		ID:          uuid.New(),
		ActionTime:  time.Now(),
		UserID:      userId,
		UserRepr:    userRepr,
		ActionFlag:  action,
		ContentType: contentType,
		ObjectID:    objectID,
		ObjectRepr:  objectRepr,
		Message:     message,
	}

	err = c.LogStore.InsertLogEntry(&logEntry)

	return nil
}

// GetPrefix returns the URL prefix for the admin panel.
func (c *AdminConfig) GetPrefix() string {
	if c.Prefix == "" {
		return ""
	}
	return "/" + c.Prefix
}

// GetAssetsPrefix returns the URL prefix for admin panel assets.
func (c *AdminConfig) GetAssetsPrefix() string {
	if c.AssetsPrefix == "" {
		return ""
	}
	return "/" + c.AssetsPrefix
}

// GetLink constructs a full link by combining the group prefix, admin prefix, and the provided link.
func (c *AdminConfig) GetLink(link string) string {
	return c.GroupPrefix + c.GetPrefix() + link
}

// GetAssetLink constructs a full asset link by combining the group prefix, assets prefix, and the file name.
func (c *AdminConfig) GetAssetLink(fileName string) string {
	return c.GroupPrefix + c.GetAssetsPrefix() + "/" + fileName
}

// GetNavBarItems generates the navigation bar items using the registered generators.
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

package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/logging"
	"github.com/google/uuid"
	"time"
)

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

type UserFetchFunction = func(ctx interface{}) (userID interface{}, repr string, err error)

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
		LogStore:                logging.NewInMemoryLogStore(100),
		LogStoreLevel:           logging.LogStoreLevelPanelView,
		NavBarGenerators:        navBarGens,
	}
}

func (c *AdminConfig) GetLogEntries(maxCount uint) []*logging.LogEntry {
	if c.LogStore == nil {
		return []*logging.LogEntry{}
	}
	entries, err := c.LogStore.GetLogEntries()
	if err != nil {
		return []*logging.LogEntry{}
	}
	entries = entries[:min(uint(len(entries)), maxCount)]
	return entries
}

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

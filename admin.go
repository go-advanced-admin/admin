package admin

import "github.com/go-advanced-admin/admin/internal/adminpanel"

// ORMIntegrator defines the interface for ORM integrations with the admin panel.
type ORMIntegrator = adminpanel.ORMIntegrator

// WebIntegrator defines the interface for web framework integrations with the admin panel.
type WebIntegrator = adminpanel.WebIntegrator

// HandlerFunc represents a handler function used in the admin panel routes.
type HandlerFunc = adminpanel.HandlerFunc

// PermissionRequest represents a request for permission to perform an action in the admin panel.
type PermissionRequest = adminpanel.PermissionRequest

// PermissionFunc defines a function type for checking permissions in the admin panel.
type PermissionFunc = adminpanel.PermissionFunc

// Panel represents the admin panel, which manages apps, models, and permissions.
type Panel = adminpanel.AdminPanel

// App represents an application within the admin panel, grouping related models together.
type App = adminpanel.App

// Config holds configuration settings for the admin panel.
type Config = adminpanel.AdminConfig

// DefaultConfig provides default configuration settings for the admin panel.
var DefaultConfig = adminpanel.DefaultAdminConfig

// NewPanel creates a new admin panel with the given ORM integrator, web integrator, permission function, and configuration.
var NewPanel = adminpanel.NewAdminPanel

// TemplateRenderer defines the interface for rendering templates in the admin panel.
type TemplateRenderer = adminpanel.TemplateRenderer

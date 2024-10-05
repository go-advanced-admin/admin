package admin

import "github.com/go-advanced-admin/admin/internal/adminpanel"

type ORMIntegrator = adminpanel.ORMIntegrator
type WebIntegrator = adminpanel.WebIntegrator

type HandlerFunc = adminpanel.HandlerFunc

type PermissionRequest = adminpanel.PermissionRequest

type PermissionFunc = adminpanel.PermissionFunc

type Panel = adminpanel.AdminPanel

type App = adminpanel.App

type Config = adminpanel.AdminConfig

var DefaultConfig = adminpanel.DefaultAdminConfig

var NewPanel = adminpanel.NewAdminPanel

type TemplateRenderer = adminpanel.TemplateRenderer

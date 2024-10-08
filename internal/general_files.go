package internal

import "embed"

//go:embed assets/*
var AssetsFiles embed.FS

//go:embed templates/*
var TemplateFiles embed.FS

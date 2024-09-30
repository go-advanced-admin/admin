package internal

import "embed"

//go:embed assets/*
var AssetsFiles embed.FS

//go:embed templates/*
var TemplateFiles embed.FS

func GetAssetsFile(fileName string) ([]byte, error) {
	return AssetsFiles.ReadFile(fileName)
}

func GetTemplateFile(fileName string) ([]byte, error) {
	return TemplateFiles.ReadFile(fileName)
}

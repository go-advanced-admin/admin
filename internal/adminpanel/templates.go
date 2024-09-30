package adminpanel

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/utils"
	"html/template"
)

type TemplateRenderer interface {
	RenderTemplate(name string, data map[string]interface{}) (string, error)
	RegisterDefaultTemplates(templates embed.FS)
	RegisterDefaultData(data map[string]interface{})
	AddCustomTemplate(name string, tmplText string) error
	RegisterDefaultAssets(assets embed.FS)
	AddCustomAsset(name string, asset []byte)
	GetAsset(name string) ([]byte, error)
	RegisterLinkFunc(func(string) string)
	RegisterAssetsFunc(func(string) string)
}

type DefaultTemplateRenderer struct {
	templates        map[string]*template.Template
	defaultTemplates embed.FS
	defaultData      map[string]interface{}
	assets           map[string]*[]byte
	defaultAssets    embed.FS
	linkFunc         func(string) string
	assetsFunc       func(string) string
}

func NewDefaultTemplateRenderer() *DefaultTemplateRenderer {
	return &DefaultTemplateRenderer{
		templates:   make(map[string]*template.Template),
		defaultData: make(map[string]interface{}),
		assets:      make(map[string]*[]byte),
	}
}

func (tr *DefaultTemplateRenderer) RenderTemplate(name string, data map[string]interface{}) (string, error) {
	newDataMap := make(map[string]interface{})
	for key, value := range tr.defaultData {
		newDataMap[key] = value
	}
	for key, value := range data {
		newDataMap[key] = value
	}
	tmpl, exists := tr.templates[name]
	if !exists {
		tmplBytes, err := tr.defaultTemplates.ReadFile(fmt.Sprintf("templates/%s", name))
		if err != nil {
			return "", err
		}
		tmpl, err = template.New(name).Funcs(tr.templateFuncs()).Parse(string(tmplBytes))
		if err != nil {
			return "", err
		}
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, newDataMap)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (tr *DefaultTemplateRenderer) RegisterDefaultTemplates(templates embed.FS) {
	tr.defaultTemplates = templates
}

func (tr *DefaultTemplateRenderer) RegisterDefaultData(data map[string]interface{}) {
	for key, value := range data {
		tr.defaultData[key] = value
	}
}

func (tr *DefaultTemplateRenderer) AddCustomTemplate(name string, tmplText string) error {
	tmpl, err := template.New(name).Funcs(tr.templateFuncs()).Parse(tmplText)
	if err != nil {
		return err
	}
	tr.templates[name] = tmpl
	return nil
}

func (tr *DefaultTemplateRenderer) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"assetPath": func(fileName string) string {
			if path, exists := tr.assets[fileName]; exists {
				return tr.assetsFunc(string(*path))
			}
			return tr.assetsFunc(fileName)
		},
		"getFieldValue": func(instance interface{}, fieldName string) (interface{}, error) {
			value, err := utils.GetFieldValue(instance, fieldName)
			if err != nil {
				return nil, err
			}
			return value, nil
		},
	}
}

func (tr *DefaultTemplateRenderer) RegisterDefaultAssets(assets embed.FS) {
	tr.defaultAssets = assets
}

func (tr *DefaultTemplateRenderer) AddCustomAsset(name string, asset []byte) {
	assetBytes := make([]byte, len(asset))
	copy(assetBytes, asset)
	tr.assets[name] = &assetBytes
}

func (tr *DefaultTemplateRenderer) GetAsset(name string) ([]byte, error) {
	assetPts, exists := tr.assets[name]
	if exists {
		return *assetPts, nil
	}
	return tr.defaultAssets.ReadFile(fmt.Sprintf("assets/%s", name))
}

func (tr *DefaultTemplateRenderer) RegisterLinkFunc(linkFunc func(string) string) {
	tr.linkFunc = linkFunc
}

func (tr *DefaultTemplateRenderer) RegisterAssetsFunc(assetsFunc func(string) string) {
	tr.assetsFunc = assetsFunc
}

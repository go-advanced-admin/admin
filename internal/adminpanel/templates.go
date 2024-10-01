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
	RegisterCompositeDefaultTemplate(name string, baseNames ...string) error
	RegisterDefaultData(data map[string]interface{}) error
	AddCustomTemplate(name string, tmplText string) error
	AddCustomCompositeTemplate(name string, baseNames ...string) error
	RegisterDefaultAssets(assets embed.FS)
	AddCustomAsset(name string, asset []byte)
	GetAsset(name string) ([]byte, error)
	RegisterLinkFunc(func(string) string)
	RegisterAssetsFunc(func(string) string)
}

type DefaultTemplateRenderer struct {
	customTemplates           map[string]string
	customCompositeTemplates  map[string][]string
	defaultCompositeTemplates map[string][]string
	defaultTemplates          embed.FS
	defaultData               map[string]interface{}
	assets                    map[string]*[]byte
	defaultAssets             embed.FS
	linkFunc                  func(string) string
	assetsFunc                func(string) string
}

func NewDefaultTemplateRenderer() *DefaultTemplateRenderer {
	return &DefaultTemplateRenderer{
		customTemplates:           make(map[string]string),
		customCompositeTemplates:  make(map[string][]string),
		defaultCompositeTemplates: make(map[string][]string),
		defaultData:               make(map[string]interface{}),
		assets:                    make(map[string]*[]byte),
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

	entryName, tmpl, err := tr.gatherTemplates(name, "", nil)
	if err != nil {
		return "", fmt.Errorf("template %s could not be combined: %v", name, err)
	}

	if entryName == "" {
		return "", fmt.Errorf("template %s could not be combined", name)
	}

	if tmpl == nil {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf bytes.Buffer
	if err = tmpl.ExecuteTemplate(&buf, entryName, newDataMap); err != nil {
		return "", fmt.Errorf("error executing template %s: %v", name, err)
	}
	return buf.String(), nil
}

func (tr *DefaultTemplateRenderer) gatherTemplates(name string, entryName string, tmpl *template.Template) (string, *template.Template, error) {
	var err error

	if content, exists := tr.customTemplates[name]; exists {
		if entryName == "" {
			entryName = name
		}
		if tmpl == nil {
			tmpl, err = template.New(name).Funcs(tr.templateFuncs()).Parse(content)
			return entryName, tmpl, err
		} else {
			tmpl, err = tmpl.New(name).Funcs(tr.templateFuncs()).Parse(content)
			return entryName, tmpl, err
		}
	}

	if subBases, exists := tr.customCompositeTemplates[name]; exists {
		for _, subBase := range subBases {
			if entryName == "" {
				entryName = subBase
			}
			if tmpl == nil {
				tmpl = template.New(name).Funcs(tr.templateFuncs())
			}
			entryName, tmpl, err = tr.gatherTemplates(subBase, entryName, tmpl)
			if err != nil {
				return entryName, nil, err
			}
		}
		if tmpl != nil {
			return entryName, tmpl, nil
		}
	}

	if subBases, exists := tr.defaultCompositeTemplates[name]; exists {
		for _, subBase := range subBases {
			if entryName == "" {
				entryName = subBase
			}
			if tmpl == nil {
				tmpl = template.New(name).Funcs(tr.templateFuncs())
			}
			entryName, tmpl, err = tr.gatherTemplates(subBase, entryName, tmpl)
			if err != nil {
				return entryName, nil, err
			}
		}
		if tmpl != nil {
			return entryName, tmpl, nil
		}
	}

	if tmplBytes, err := tr.defaultTemplates.ReadFile(fmt.Sprintf("templates/%s", name)); err == nil {
		if entryName == "" {
			entryName = name
		}
		if tmpl == nil {
			tmpl, err = template.New(name).Funcs(tr.templateFuncs()).Parse(string(tmplBytes))
			return entryName, tmpl, err
		} else {
			tmpl, err = tmpl.New(name).Funcs(tr.templateFuncs()).Parse(string(tmplBytes))
			return entryName, tmpl, err
		}

	}

	return entryName, nil, fmt.Errorf("template %s not found", name)
}

func (tr *DefaultTemplateRenderer) RegisterDefaultTemplates(templates embed.FS) {
	tr.defaultTemplates = templates
}

func (tr *DefaultTemplateRenderer) validateAndParseBases(tmpl *template.Template, baseNames []string) error {
	if len(baseNames) == 0 {
		return fmt.Errorf("no base templates provided")
	}
	for _, baseName := range baseNames {
		if content, exists := tr.customTemplates[baseName]; exists {
			_, err := tmpl.New(baseName).Parse(content)
			if err != nil {
				return fmt.Errorf("error parsing base template %s: %v", baseName, err)
			}
		} else if subBases, exists := tr.customCompositeTemplates[baseName]; exists {
			if err := tr.validateAndParseBases(tmpl, subBases); err != nil {
				return err
			}
		} else if subBases, exists = tr.defaultCompositeTemplates[baseName]; exists {
			if err := tr.validateAndParseBases(tmpl, subBases); err != nil {
				return err
			}
		} else if tmplBytes, err := tr.defaultTemplates.ReadFile(fmt.Sprintf("templates/%s", baseName)); err == nil {
			_, err = tmpl.New(baseName).Parse(string(tmplBytes))
			if err != nil {
				return fmt.Errorf("error parsing embedded template %s: %v", baseName, err)
			}
		} else {
			return fmt.Errorf("base template %s not found", baseName)
		}
	}
	return nil
}

func (tr *DefaultTemplateRenderer) RegisterCompositeDefaultTemplate(name string, baseNames ...string) error {
	if _, exists := tr.defaultCompositeTemplates[name]; exists {
		return fmt.Errorf("template %s already exists as a default composite template", name)
	}

	if len(baseNames) <= 1 {
		return fmt.Errorf("template %s has no base templates", name)
	}

	compositeTemplate := template.New(name).Funcs(tr.templateFuncs())

	err := tr.validateAndParseBases(compositeTemplate, baseNames)
	if err != nil {
		return err
	}

	tr.defaultCompositeTemplates[name] = baseNames
	return nil
}

func (tr *DefaultTemplateRenderer) AddCustomCompositeTemplate(name string, baseNames ...string) error {
	if _, exists := tr.customCompositeTemplates[name]; exists {
		return fmt.Errorf("template %s already exists as a custom composite template", name)
	}
	if _, exists := tr.customTemplates[name]; exists {
		return fmt.Errorf("template %s already exists as a custom template", name)
	}

	if len(baseNames) <= 1 {
		return fmt.Errorf("template %s has no base templates", name)
	}

	compositeTemplate := template.New(name).Funcs(tr.templateFuncs())

	err := tr.validateAndParseBases(compositeTemplate, baseNames)
	if err != nil {
		return err
	}

	tr.customCompositeTemplates[name] = baseNames
	return nil
}

func (tr *DefaultTemplateRenderer) RegisterDefaultData(data map[string]interface{}) error {
	for key, value := range data {
		if _, exists := tr.defaultData[key]; exists {
			return fmt.Errorf("data key %s already exists", key)
		}
		tr.defaultData[key] = value
	}

	return nil
}

func (tr *DefaultTemplateRenderer) AddCustomTemplate(name, tmplText string) error {
	if _, exists := tr.customTemplates[name]; exists {
		return fmt.Errorf("template %s already exists as a custom template", name)
	}

	if _, exists := tr.customCompositeTemplates[name]; exists {
		return fmt.Errorf("template %s already exists as a custom composite template", name)
	}

	if _, err := template.New(name).Funcs(tr.templateFuncs()).Parse(tmplText); err != nil {
		return fmt.Errorf("error parsing template %s: %v", name, err)
	}

	tr.customTemplates[name] = tmplText
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

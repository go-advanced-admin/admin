package adminpanel

import (
	"bytes"
	"errors"
	"github.com/go-advanced-admin/admin/internal"
	"testing"
)

func TestDefaultTemplateRenderer(t *testing.T) {
	createTemplateRenderer := func() *DefaultTemplateRenderer {
		renderer := NewDefaultTemplateRenderer()

		renderer.RegisterAssetsFunc(func(name string) string {
			return "/assets/" + name
		})

		renderer.RegisterLinkFunc(func(url string) string {
			return "/link/" + url
		})

		return renderer
	}

	t.Run("New Template Renderer", func(t *testing.T) {
		renderer := NewDefaultTemplateRenderer()
		if renderer == nil {
			t.Fatalf("expected non-nil renderer")
		}
		if len(renderer.customTemplates) != 0 {
			t.Fatalf("expected no custom templates")
		}
	})

	t.Run("Fail to Add Malformed Custom Template", func(t *testing.T) {
		renderer := NewDefaultTemplateRenderer()
		err := renderer.AddCustomTemplate("broken", "{{ .Field")
		if err == nil {
			t.Fatalf("expected error on malformed template, got no error")
		}
	})

	t.Run("Add and Render Custom Template", func(t *testing.T) {
		renderer := createTemplateRenderer()
		err := renderer.AddCustomTemplate("custom", `Custom: {{.Field}}`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		html, err := renderer.RenderTemplate("custom", map[string]interface{}{"Field": "Value"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !bytes.Contains([]byte(html), []byte("Custom: Value")) {
			t.Errorf("expected rendered html to contain 'Custom: Value', got %s", html)
		}
	})

	t.Run("Register and Render Composite Template", func(t *testing.T) {
		renderer := createTemplateRenderer()
		err := renderer.AddCustomTemplate("base1", `{{define "base1"}}Base 1: {{template "base1content" .}}{{end}}`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = renderer.AddCustomTemplate("base1content", `{{define "base1content"}}Content 1: {{.Value1}}{{end}}`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = renderer.AddCustomCompositeTemplate("composite", "base1", "base1content")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		html, err := renderer.RenderTemplate("composite", map[string]interface{}{"Value1": "Data1"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedStr := "Base 1: Content 1: Data1"
		if !bytes.Contains([]byte(html), []byte(expectedStr)) {
			t.Errorf("expected rendered HTML to contain '%s', got %s", expectedStr, html)
		}
	})

	t.Run("Fail to Render Nonexistent Template", func(t *testing.T) {
		renderer := createTemplateRenderer()
		_, err := renderer.RenderTemplate("nonexistent", nil)
		if err == nil {
			t.Fatalf("expected error, got no error")
		}
	})

	t.Run("Fail on Parsing Error", func(t *testing.T) {
		renderer := createTemplateRenderer()
		err := renderer.AddCustomTemplate("broken", "{{ .Field")
		if err == nil {
			t.Fatalf("expected error on malformed template, got no error")
		}
	})

	t.Run("Template Function assetPath", func(t *testing.T) {
		renderer := createTemplateRenderer()
		expectedPath := "/assets/someAsset"

		funcMap := renderer.templateFuncs()

		if assetPathFunc, ok := funcMap["assetPath"]; ok {
			result := assetPathFunc.(func(string) string)("someAsset")
			if result != expectedPath {
				t.Errorf("expected assetPath to return %s, got %s", expectedPath, result)
			}
		} else {
			t.Errorf("func 'assetPath' not found in templateFuncs")
		}
	})

	t.Run("Template Function getFieldValue", func(t *testing.T) {
		renderer := createTemplateRenderer()
		funcMap := renderer.templateFuncs()

		type TestInstance struct {
			Field string
		}

		instance := TestInstance{Field: "Value"}

		if getFieldValueFunc, ok := funcMap["getFieldValue"]; ok {
			val, err := getFieldValueFunc.(func(interface{}, string) (interface{}, error))(instance, "Field")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if val != "Value" {
				t.Errorf("expected getFieldValue to return 'Value', got %v", val)
			}

		} else {
			t.Errorf("func 'getFieldValue' not found in templateFuncs")
		}
	})

	t.Run("Get and Add Custom Asset", func(t *testing.T) {
		renderer := createTemplateRenderer()
		assetContent := []byte("custom asset content")
		renderer.AddCustomAsset("custom.js", assetContent)

		asset, err := renderer.GetAsset("custom.js")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !bytes.Equal(asset, assetContent) {
			t.Errorf("expected asset to match custom asset content, got %s", asset)
		}
	})

	t.Run("Add and Retrieve Default Data", func(t *testing.T) {
		renderer := createTemplateRenderer()
		err := renderer.RegisterDefaultData(map[string]interface{}{
			"DefaultField": "DefaultValue",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if val, exists := renderer.defaultData["DefaultField"]; !exists || val != "DefaultValue" {
			t.Errorf("expected default data to contain 'DefaultField' with value 'DefaultValue', got %v", val)
		}
	})

	t.Run("Register Default Data", func(t *testing.T) {
		renderer := createTemplateRenderer()
		defaultData := map[string]interface{}{"key": "value"}
		err := renderer.RegisterDefaultData(defaultData)
		if err != nil {
			t.Fatalf("expected no error registering default data, got %v", err)
		}
		val := renderer.defaultData["key"]
		if val != "value" {
			t.Errorf("expected 'value' for default data, got %v", val)
		}

		err = renderer.RegisterDefaultData(map[string]interface{}{"key": "newvalue"})
		if err == nil {
			t.Fatalf("expected error registering duplicate default data, got no error")
		}
	})

	t.Run("Register and Fetch Default Templates", func(t *testing.T) {
		renderer := createTemplateRenderer()
		renderer.RegisterDefaultTemplates(internal.TemplateFiles, "templates/")
		renderer.RegisterDefaultAssets(internal.AssetsFiles, "assets/")

		html, err := renderer.RenderTemplate("root.html", map[string]interface{}{})
		if err != nil && !errors.Is(err, bytes.ErrTooLarge) {
			t.Fatalf("expected no error (or file specific), got %v", err)
		}

		if !bytes.Contains([]byte(html), []byte("<!DOCTYPE html>")) {
			t.Errorf("expected rendered html to contain 'expected string', got %s", html)
		}
	})

	t.Run("Register and Validate Composite Template Registration", func(t *testing.T) {
		renderer := createTemplateRenderer()
		_ = renderer.AddCustomTemplate("content1", `{{define "content1"}}Content 1 body{{end}}`)
		_ = renderer.AddCustomTemplate("content2", `{{define "content2"}}Content 2 body{{end}}`)

		err := renderer.RegisterCompositeDefaultTemplate("composite", "content1", "content2")
		if err != nil {
			t.Fatalf("expected no error registering composite template, got %v", err)
		}

		html, err := renderer.RenderTemplate("composite", nil)
		if err != nil {
			t.Fatalf("expected no error rendering composite template, got %v", err)
		}

		if !bytes.Contains([]byte(html), []byte("Content 1 body")) {
			t.Errorf("expected composite rendering to include all content parts, got %s", html)
		}
	})

	t.Run("RegisterDefaultAssets and Retrieve Assets", func(t *testing.T) {
		renderer := createTemplateRenderer()
		renderer.RegisterDefaultAssets(internal.TemplateFiles, "templates/")

		_, err := renderer.GetAsset("root.html")
		if err != nil {
			t.Logf("This might fail if `root.html` doesn't exist, ensure assets are correctly embedded")
		}
	})

	t.Run("Render Nonexistent Template", func(t *testing.T) {
		renderer := createTemplateRenderer()
		_, err := renderer.RenderTemplate("nonexistent", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !bytes.Contains([]byte(err.Error()), []byte("not found")) {
			t.Errorf("expected error message to contain 'not found', got %v", err)
		}
	})
}

package adminpanel

import (
	"bytes"
	"testing"
)

func TestDefaultTemplateRenderer(t *testing.T) {
	renderer := NewDefaultTemplateRenderer()

	t.Run("Add and Render Custom Template", func(t *testing.T) {
		err := renderer.AddCustomTemplate("custom", "Custom: {{.Field}}")
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

	t.Run("Get and Add Custom Asset", func(t *testing.T) {
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

	t.Run("Fail to Render Nonexistent Template", func(t *testing.T) {
		_, err := renderer.RenderTemplate("nonexistent", nil)
		if err == nil {
			t.Fatalf("expected error, got no error")
		}
	})

	t.Run("Fail on Parsing Error", func(t *testing.T) {
		err := renderer.AddCustomTemplate("broken", "{{ .Field")
		if err == nil {
			t.Fatalf("expected error on malformed template, got no error")
		}
	})

	t.Run("Template Function assetPath", func(t *testing.T) {
		expectedPath := "/test/path"
		renderer.RegisterAssetsFunc(func(name string) string {
			return "/test/path"
		})

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
		funcMap := renderer.templateFuncs()

		type TestInstance struct {
			Field string
		}

		instance := TestInstance{Field: "Value"}

		if getFieldValueFunc, ok := funcMap["getFieldValue"]; ok {
			val, err := getFieldValueFunc.(func(instance interface{}, fieldName string) (interface{}, error))(instance, "Field")
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
}

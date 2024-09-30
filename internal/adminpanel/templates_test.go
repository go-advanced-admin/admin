package adminpanel

import (
	"github.com/go-advanced-admin/admin/internal"
	"testing"
)

func TestDefaultTemplateRenderer_RenderTemplate(t *testing.T) {
	renderer := NewDefaultTemplateRenderer()
	renderer.RegisterDefaultTemplates(internal.TemplateFiles)
	renderer.RegisterDefaultAssets(internal.AssetsFiles)
	renderer.RegisterAssetsFunc(func(input string) string {
		return input
	})

	html, err := renderer.RenderTemplate("root.html", map[string]interface{}{"Title": "Test"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if html == "" {
		t.Errorf("expected HTML string, got empty string")
	}
}

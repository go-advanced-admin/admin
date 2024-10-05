package adminpanel

import (
	"testing"
)

func TestAdminConfig_GetPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		expected string
	}{
		{"Non-empty Prefix", "admin", "/admin"},
		{"Empty Prefix", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := AdminConfig{Prefix: tt.prefix}
			if got := config.GetPrefix(); got != tt.expected {
				t.Errorf("GetPrefix() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAdminConfig_GetAssetsPrefix(t *testing.T) {
	tests := []struct {
		name         string
		assetsPrefix string
		expected     string
	}{
		{"Non-empty AssetsPrefix", "assets", "/assets"},
		{"Empty AssetsPrefix", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := AdminConfig{AssetsPrefix: tt.assetsPrefix}
			if got := config.GetAssetsPrefix(); got != tt.expected {
				t.Errorf("GetAssetsPrefix() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAdminConfig_GetLink(t *testing.T) {
	tests := []struct {
		name        string
		groupPrefix string
		prefix      string
		link        string
		expected    string
	}{
		{"With Prefix and Group Prefix", "group", "admin", "/dashboard", "group/admin/dashboard"},
		{"Without Prefix", "group", "", "/dashboard", "group/dashboard"},
		{"Without Group Prefix", "", "admin", "/dashboard", "/admin/dashboard"},
		{"Without Both Prefixes", "", "", "/dashboard", "/dashboard"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := AdminConfig{GroupPrefix: tt.groupPrefix, Prefix: tt.prefix}
			if got := config.GetLink(tt.link); got != tt.expected {
				t.Errorf("GetLink() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAdminConfig_GetAssetLink(t *testing.T) {
	tests := []struct {
		name         string
		groupPrefix  string
		assetsPrefix string
		fileName     string
		expected     string
	}{
		{"With AssetsPrefix and Group Prefix", "group", "assets", "style.css", "group/assets/style.css"},
		{"Without AssetsPrefix", "group", "", "style.css", "group/style.css"},
		{"Without Group Prefix", "", "assets", "style.css", "/assets/style.css"},
		{"Without Both Prefixes", "", "", "style.css", "/style.css"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := AdminConfig{GroupPrefix: tt.groupPrefix, AssetsPrefix: tt.assetsPrefix}
			if got := config.GetAssetLink(tt.fileName); got != tt.expected {
				t.Errorf("GetAssetLink() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package adminpanel

import "testing"

func TestAdminConfig_GetPrefix(t *testing.T) {
	config := AdminConfig{Prefix: "admin"}
	if config.GetPrefix() != "/admin" {
		t.Errorf("expected /admin, got %s", config.GetPrefix())
	}

	config = AdminConfig{Prefix: ""}
	if config.GetPrefix() != "" {
		t.Errorf("expected empty string, got %s", config.GetPrefix())
	}
}

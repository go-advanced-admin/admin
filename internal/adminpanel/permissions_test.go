package adminpanel

import "testing"

func TestPermissionFunc_HasPermission(t *testing.T) {
	permFunc := PermissionFunc(func(req PermissionRequest, ctx interface{}) (bool, error) {
		return req.Action != nil && *req.Action == ReadAction, nil
	})

	action := ReadAction
	req := PermissionRequest{Action: &action}
	allowed, err := permFunc.HasPermission(req, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !allowed {
		t.Fatalf("expected permission to be allowed")
	}

	req = PermissionRequest{Action: nil}
	allowed, err = permFunc.HasPermission(req, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if allowed {
		t.Fatalf("expected permission to be denied")
	}
}

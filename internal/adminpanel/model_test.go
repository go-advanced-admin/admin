package adminpanel

import "testing"

type TestModel struct {
	ID   uint
	Name string
}

func TestGetPrimaryKeyGetter(t *testing.T) {
	model := &TestModel{}
	getter, err := GetPrimaryKeyGetter(model)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	id := getter(model)
	if id != uint(0) {
		t.Errorf("expected 0, got %v", id)
	}
}

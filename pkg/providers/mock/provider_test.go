package mock

import "testing"

func TestAuthTrue(t *testing.T) {
	provider := NewMockProvider(true)
	if provider.CheckAuth() != nil {
		t.Error("unexpected error")
	}
}

func TestAuthFalse(t *testing.T) {
	provider := NewMockProvider(false)
	if provider.CheckAuth() == nil {
		t.Error("unexpected success")
	}
}

func TestCreate(t *testing.T) {
	provider := NewMockProvider(false)
	new := MockRecord{
		Name: "asdf",
		Type: "A",
	}

	provider.CreateRecord("", new)

	records, err := provider.AllRecords("")
	if err != nil {
		t.Error("unexpected error")
	}
	if len(records) != 1 {
		t.Error("unexpected length", len(records))
	}
}

func TestDelete(t *testing.T) {
	provider := NewMockProvider(false)
	new := MockRecord{
		Name: "asdf",
		Type: "A",
	}

	provider.CreateRecord("", new)
	provider.DeleteRecord("", new)

	records, err := provider.AllRecords("")
	if err != nil {
		t.Error("unexpected error")
	}
	if len(records) != 0 {
		t.Error("unexpected length", len(records))
	}
}

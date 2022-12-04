package secrets

import "testing"

const (
	mockKey = "key_1"
	mockVal = "val_1"
)

func TestProvider(t *testing.T) {
	provider := Provider{}

	provider.Set(mockKey, mockKey)

	if provider.Read(mockKey) != mockKey {
		t.Error("value not found")
	}
}

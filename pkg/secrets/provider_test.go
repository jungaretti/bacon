package secrets

import "testing"

const (
	mockKey = "key_1"
	mockVal = "val_1"
)

func TestProvider(t *testing.T) {
	provider := Provider{}

	provider.Set(mockKey, mockVal)

	if provider.Read(mockKey) != mockVal {
		t.Error("value not found")
	}
}

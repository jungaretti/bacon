package cmd

import (
	"bacon/pkg/providers/mock"
	"testing"
)

func TestPrint(t *testing.T) {
	mockProvider := mock.NewMockProvider(true)
	mockApp := App{
		Provider: mockProvider,
	}

	if Print(&mockApp, "") != nil {
		t.Error("unexected error")
	}
}

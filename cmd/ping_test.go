package cmd

import (
	"bacon/pkg/providers/mock"
	"testing"
)

func TestPingSuccess(t *testing.T) {
	mockProvider := mock.NewMockProvider(true)
	mockApp := App{
		Provider: mockProvider,
	}

	if Ping(&mockApp) != nil {
		t.Error("unexpected error")
	}
}

func TestPingFailure(t *testing.T) {
	mockProvider := mock.NewMockProvider(false)
	mockApp := App{
		Provider: mockProvider,
	}

	if Ping(&mockApp) == nil {
		t.Error("unexpected success")
	}
}

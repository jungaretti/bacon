package porkbun

import (
	"testing"
)

const (
	mockType    = "mock_type"
	mockContent = "mock_content"
	mockTtl     = "mock_ttl"
	mockName    = "mock_host"
)

func TestRecord(t *testing.T) {
	porkRecord := Record{
		Type:    mockType,
		Content: mockContent,
		TTL:     mockTtl,
		Name:    mockName,
	}

	if porkRecord.GetType() != mockType {
		t.Error("did not return the expected type")
	}
	if porkRecord.GetContent() != mockContent {
		t.Error("did not return the expected content")
	}
	if porkRecord.GetTTL() != mockTtl {
		t.Error("did not return the expected TTL")
	}
	if porkRecord.GetHost() != mockName {
		t.Error("did not return the expected host")
	}
}

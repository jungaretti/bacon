package record

import (
	"testing"
)

const (
	mockName    = "mock_host"
	mockType    = "mock_type"
	mockTtl     = "mock_ttl"
	mockContent = "mock_content"
)

func TestRecord(t *testing.T) {
	porkRecord := Record{
		Name:    mockName,
		Type:    mockType,
		TTL:     mockTtl,
		Content: mockContent,
	}

	if porkRecord.GetName() != mockName {
		t.Error("did not return the expected host")
	}
	if porkRecord.GetType() != mockType {
		t.Error("did not return the expected type")
	}
	if porkRecord.GetTtl() != mockTtl {
		t.Error("did not return the expected TTL")
	}
	if porkRecord.GetData() != mockContent {
		t.Error("did not return the expected content")
	}
}

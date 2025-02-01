package api

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

func TestIgnoredName(t *testing.T) {
	ignoredRecord := Record{
		Name:    "_acme-challenge.bacontest42.com",
		Type:    "TXT",
		TTL:     "600",
		Content: "c_V4WaKPWlisAvnvTZ62BOuLiQDpkC2cOtahW8TDthw",
	}

	if !ignoredRecord.isIgnored() {
		t.Error("did not ignore record with _acme-challenge")
	}
}

func TestIgnoredType(t *testing.T) {
	ignoredRecord := Record{
		Name:    "www.bacontest42.com",
		Type:    "NS",
		TTL:     "600",
		Content: "c_V4WaKPWlisAvnvTZ62BOuLiQDpkC2cOtahW8TDsfs",
	}

	if !ignoredRecord.isIgnored() {
		t.Error("did not ignore NS record")
	}
}

package porkbun

import (
	"testing"
)

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

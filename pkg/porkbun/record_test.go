package porkbun

import (
	"testing"
)

func TestRecordHashDiffsName(t *testing.T) {
	recordA := Record{
		Name:    "www.bacontest42.com",
		Type:    "A",
		TTL:     "600",
		Content: "1.2.3.4",
	}

	recordB := Record{
		Name:    "mail.bacontest42.com",
		Type:    recordA.Type,
		TTL:     recordA.TTL,
		Content: recordA.Content,
	}

	if RecordHash(recordA) == RecordHash(recordB) {
		t.Error("records that differ only by name have the same hash")
	}
}

func TestRecordHashDiffsType(t *testing.T) {
	recordA := Record{
		Name:    "www.bacontest42.com",
		Type:    "A",
		TTL:     "600",
		Content: "1.2.3.4",
	}

	recordB := Record{
		Name:    recordA.Name,
		Type:    "CNAME",
		TTL:     recordA.TTL,
		Content: "somewhere.else.com",
	}

	if RecordHash(recordA) == RecordHash(recordB) {
		t.Error("records that differ only by type have the same hash")
	}
}

func TestRecordHashDiffsContent(t *testing.T) {
	recordA := Record{
		Name:    "www.bacontest42.com",
		Type:    "A",
		TTL:     "600",
		Content: "1.2.3.4",
	}

	recordB := Record{
		Name:    recordA.Name,
		Type:    recordA.Type,
		TTL:     recordA.TTL,
		Content: "5.6.7.8",
	}

	if RecordHash(recordA) == RecordHash(recordB) {
		t.Error("records that differ only by content have the same hash")
	}
}

func TestRecordHashDiffsTTL(t *testing.T) {
	recordA := Record{
		Name:    "www.bacontest42.com",
		Type:    "A",
		TTL:     "600",
		Content: "1.2.3.4",
	}

	recordB := Record{
		Name:    recordA.Name,
		Type:    recordA.Type,
		TTL:     "3600",
		Content: recordA.Content,
	}

	if RecordHash(recordA) == RecordHash(recordB) {
		t.Error("records that differ only by TTL have the same hash")
	}
}

func TestRecordHashDiffsPriority(t *testing.T) {
	recordA := Record{
		Name:     "mail.bacontest42.com",
		Type:     "MX",
		TTL:      "600",
		Content:  "mx.example.com",
		Priority: "10",
	}

	recordB := Record{
		Name:     recordA.Name,
		Type:     recordA.Type,
		TTL:      recordA.TTL,
		Content:  recordA.Content,
		Priority: "20",
	}

	if RecordHash(recordA) == RecordHash(recordB) {
		t.Error("records that differ only by priority have the same hash")
	}
}

func TestRecordHashDiffsNotes(t *testing.T) {
	recordA := Record{
		Name:     "mail.bacontest42.com",
		Type:     "MX",
		TTL:      "600",
		Content:  "mx.example.com",
		Priority: "10",
		Notes:    "Fastmail",
	}

	recordB := Record{
		Name:     recordA.Name,
		Type:     recordA.Type,
		TTL:      recordA.TTL,
		Content:  recordA.Content,
		Priority: recordA.Priority,
		Notes:    "Exchange",
	}

	if RecordHash(recordA) == RecordHash(recordB) {
		t.Error("records that differ only by notes have the same hash")
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

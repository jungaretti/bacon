package porkbun

import (
	"testing"
)

func TestRecordHashNoDiffs(t *testing.T) {
	recordA := Record{
		Name:    "www.bacondemo.com",
		Type:    "A",
		TTL:     "600",
		Content: "1.2.3.4",
	}

	recordB := Record{
		Name:    recordA.Name,
		Type:    recordA.Type,
		TTL:     recordA.TTL,
		Content: recordA.Content,
	}

	if RecordHash(recordA) != RecordHash(recordB) {
		t.Error("records that are the same have different hashes")
	}
}

func TestRecordHashDiffsName(t *testing.T) {
	recordA := Record{
		Name:    "www.bacondemo.com",
		Type:    "A",
		TTL:     "600",
		Content: "1.2.3.4",
	}

	recordB := Record{
		Name:    "mail.bacondemo.com",
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
		Name:    "www.bacondemo.com",
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
		Name:    "www.bacondemo.com",
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
		Name:    "www.bacondemo.com",
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
		Name:     "mail.bacondemo.com",
		Type:     "MX",
		TTL:      "600",
		Content:  "mx.bacondemo.com",
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
		Name:     "mail.bacondemo.com",
		Type:     "MX",
		TTL:      "600",
		Content:  "mx.bacondemo.com",
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

func TestRecordIdentityHashNoDiffs(t *testing.T) {
	recordA := Record{
		Name: "www.bacondemo.com",
		Type: "A",
	}

	recordB := Record{
		Name: recordA.Name,
		Type: recordA.Type,
	}

	if RecordIdentityHash(recordA) != RecordIdentityHash(recordB) {
		t.Error("records that are the same have different identity hashes")
	}
}

func TestRecordIdentityHashDiffsName(t *testing.T) {
	recordA := Record{
		Name: "www.bacondemo.com",
		Type: "A",
	}

	recordB := Record{
		Name: "mail.bacondemo.com",
		Type: recordA.Type,
	}

	if RecordIdentityHash(recordA) == RecordIdentityHash(recordB) {
		t.Error("records that differ by name have the same identity hash")
	}
}

func TestRecordIdentityHashDiffsType(t *testing.T) {
	recordA := Record{
		Name: "www.bacondemo.com",
		Type: "A",
	}

	recordB := Record{
		Name: recordA.Name,
		Type: "CNAME",
	}

	if RecordIdentityHash(recordA) == RecordIdentityHash(recordB) {
		t.Error("records that differ by type have the same identity hash")
	}
}

func TestIgnoredName(t *testing.T) {
	ignoredRecord := Record{
		Name:    "_acme-challenge.bacondemo.com",
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
		Name:    "www.bacondemo.com",
		Type:    "NS",
		TTL:     "600",
		Content: "c_V4WaKPWlisAvnvTZ62BOuLiQDpkC2cOtahW8TDsfs",
	}

	if !ignoredRecord.isIgnored() {
		t.Error("did not ignore NS record")
	}
}

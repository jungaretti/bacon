package config

import (
	"bacon/pkg/porkbun"
	"testing"
)

func TestRecordToPorkbun(t *testing.T) {
	record := Record{
		Name: "bacontest42.com",
		Type: "A",
		Ttl:  600,
		Data: "1.2.3.4",
	}

	porkRecord := record.ToPorkbun()
	if porkRecord.Name != record.Name {
		t.Error("expected name", record.Name, "got", porkRecord.Name)
	}
	if porkRecord.Type != record.Type {
		t.Error("expected type", record.Type, "got", porkRecord.Type)
	}
	if porkRecord.TTL != "600" {
		t.Error("expected TTL 600, got", porkRecord.TTL)
	}
	if porkRecord.Content != record.Data {
		t.Error("expected content", record.Data, "got", porkRecord.Content)
	}
}

func TestRecordToPorkbunWithPriority(t *testing.T) {
	record := Record{
		Name:     "bacontest42.com",
		Type:     "MX",
		Ttl:      600,
		Data:     "mail.bacontest42.com",
		Priority: 10,
	}

	porkRecord := record.ToPorkbun()
	if porkRecord.Priority != "10" {
		t.Error("expected priority 10, got", porkRecord.Priority)
	}
}

func TestRecordToPorkbunWithNotes(t *testing.T) {
	record := Record{
		Name:  "capybara.bacontest42.com",
		Type:  "A",
		Ttl:   600,
		Data:  "1.2.3.4",
		Notes: "Capybara",
	}

	porkRecord := record.ToPorkbun()
	if porkRecord.Notes != "Capybara" {
		t.Error("expected notes 'Capybara', got", porkRecord.Notes)
	}
}

func TestRecordFromPorkbun(t *testing.T) {
	record := porkbun.Record{
		Name:    "bacontest42.com",
		Type:    "A",
		TTL:     "600",
		Content: "1.2.3.4",
	}

	configRecord, err := RecordFromPorkbun(record)
	if err != nil {
		t.Fatalf("converting record from Porkbun: %v", err)
	}

	if configRecord.Name != record.Name {
		t.Error("expected name", record.Name, "got", configRecord.Name)
	}
	if configRecord.Type != record.Type {
		t.Error("expected type", record.Type, "got", configRecord.Type)
	}
	if configRecord.Ttl != 600 {
		t.Error("expected TTL 600, got", configRecord.Ttl)
	}
	if configRecord.Data != record.Content {
		t.Error("expected content", record.Content, "got", configRecord.Data)
	}
}

func TestRoundTripConversion(t *testing.T) {
	original := Record{
		Name: "bacontest42.com",
		Type: "A",
		Ttl:  600,
		Data: "1.2.3.4",
	}

	porkRecord := original.ToPorkbun()
	converted, err := RecordFromPorkbun(porkRecord)
	if err != nil {
		t.Fatalf("converting record from Porkbun: %v", err)
	}

	if original != converted {
		t.Errorf("round trip conversion failed: expected %+v, got %+v", original, converted)
	}
}

package config

import (
	"bacon/pkg/porkbun"
	"os"
	"testing"
)

func TestValidConfig(t *testing.T) {
	configFile, err := seedConfigToTempFile(`
domain: bacontest42.com
records:
    - host: bacontest42.com
      type: ALIAS
      ttl: 600
      content: pixie.porkbun.com
    - host: '*.bacontest42.com'
      type: CNAME
      ttl: 600
      content: pixie.porkbun.com
    - type: MX
      host: bacontest42.com
      content: in1-smtp.messagingengine.com
      ttl: 600
      priority: 10
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	config, err := ReadFile(configFile)
	if err != nil {
		t.Fatal("could not read config file", err)
	}

	if len(config.Records) != 3 {
		t.Fatal("expected 3 records, got", len(config.Records))
	}
}

func TestInvalidConfig(t *testing.T) {
	configFile, err := seedConfigToTempFile(`
domain: bacontest42.com
records:
    - type: ALIAS
      ttl: 600
      content: pixie.porkbun.com
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	_, err = ReadFile(configFile)
	if err == nil {
		t.Fatal("expected error when config is invalid", err)
	}
}

func TestRecordToPorkbunWithPriority(t *testing.T) {
	record := Record{
		Name:     "bacontest42.com",
		Type:     "MX",
		Ttl:      600,
		Data:     "in1-smtp.messagingengine.com",
		Priority: 10,
	}

	porkRecord := record.ToPorkbun()
	if porkRecord.Priority != "10" {
		t.Error("expected priority 10, got", porkRecord.Priority)
	}
}

func TestRecordToPorkbunWithoutPriority(t *testing.T) {
	record := Record{
		Name: "*.bacontest42.com",
		Type: "CNAME",
		Ttl:  600,
		Data: "pixie.porkbun.com",
	}

	porkRecord := record.ToPorkbun()
	if porkRecord.Priority != "" {
		t.Error("expected empty priority, got", porkRecord.Priority)
	}
}

func TestRecordFromPorkbunWithPriority(t *testing.T) {
	record := porkbun.Record{
		Name:     "bacontest42.com",
		Type:     "MX",
		TTL:      "600",
		Content:  "in1-smtp.messagingengine.com",
		Priority: "10",
	}

	configRecord, err := RecordFromPorkbun(record)
	if err != nil {
		t.Fatal("did not convert record", err)
	}
	if configRecord.Priority != 10 {
		t.Error("expected priority 10, got", configRecord.Priority)
	}
}

func TestRecordFromPorkbunWithoutPriority(t *testing.T) {
	record := porkbun.Record{
		Name:     "*.bacontest42.com",
		Type:     "CNAME",
		TTL:      "600",
		Content:  "pixie.porkbun.com",
		Priority: "0",
	}

	configRecord, err := RecordFromPorkbun(record)
	if err != nil {
		t.Fatal("did not convert record", err)
	}
	if configRecord.Priority != 0 {
		t.Error("expected priority 0, got", configRecord.Priority)
	}
}

func seedConfigToTempFile(mockConfig string) (string, error) {
	tempFile, err := os.CreateTemp("", "tmpfile-*")
	if err != nil {
		return "", err
	}

	defer tempFile.Close()

	_, err = tempFile.WriteString(mockConfig)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

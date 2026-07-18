package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidConfig(t *testing.T) {
	configFile := seedConfigToTempFile(t, `
domain: bacondemo.com
records:
    - host: bacondemo.com
      type: ALIAS
      ttl: 600
      content: pixie.porkbun.com
    - host: '*.bacondemo.com'
      type: CNAME
      ttl: 600
      content: pixie.porkbun.com
    - type: MX
      host: bacondemo.com
      content: in1-smtp.messagingengine.com
      ttl: 600
      priority: 10
`)

	config, err := ReadFile(configFile)
	if err != nil {
		t.Fatal("could not read config file", err)
	}

	if len(config.Records) != 3 {
		t.Fatal("expected 3 records, got", len(config.Records))
	}
}

func TestInvalidConfig(t *testing.T) {
	configFile := seedConfigToTempFile(t, `
domain: bacondemo.com
records:
    - type: ALIAS
      ttl: 600
      content: pixie.porkbun.com
`)

	_, err := ReadFile(configFile)
	if err == nil {
		t.Fatal("expected error when config is invalid", err)
	}
}

func seedConfigToTempFile(t *testing.T, mockConfig string) string {
	t.Helper()

	configFile := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(configFile, []byte(mockConfig), 0600); err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	return configFile
}

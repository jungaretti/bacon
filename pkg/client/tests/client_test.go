package client

import (
	"bacon/pkg/client"
	"os"
	"testing"
)

const TEST_CONFIG_FILE = `domain: example.com
records:
- type: CNAME
  host: example.com
  content: example.org
  ttl: "600"
- type: MX
  host: example.com
  content: fwd2.example.org
  ttl: "600"
  priority: "20"`

func TestReadConfig(t *testing.T) {
	temp, err := os.CreateTemp("", "baconconfig-")
	if err != nil {
		t.Fatalf("failed to create new temp file")
	}

	_, err = temp.WriteString(TEST_CONFIG_FILE)
	if err != nil {
		t.Fatalf("failed to write test config")
	}

	config, err := client.ReadConfig(temp.Name())
	if err != nil {
		t.Fatalf("failed to read config")
	}

	if config.Domain != "example.com" {
		t.Error("wrong domain", config.Domain)
	}
	if len(config.Records) != 2 {
		t.Error("wrong number of records", len(config.Records))
	}
}

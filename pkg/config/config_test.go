package config

import (
	"testing"
)

func TestValidConfig(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
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
		t.Fatal("expected 3 records after deployment, got", len(config.Records))
	}
}

func TestInvalidConfigMissingHost(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
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
		t.Fatal("expected error when a record is missing host field", err)
	}
}

func TestInvalidConfigMissingType(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
domain: bacontest42.com
records:
    - host: bacontest42.com
      ttl: 600
      content: pixie.porkbun.com
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	_, err = ReadFile(configFile)
	if err == nil {
		t.Fatal("expected error when a record is missing type field", err)
	}
}

func TestInvalidConfigInvalidType(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
domain: bacontest42.com
records:
    - host: bacontest42.com
      type: FAKE
      ttl: 600
      content: pixie.porkbun.com
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	_, err = ReadFile(configFile)
	if err == nil {
		t.Fatal("expected error when a record has an invalid type", err)
	}
}

func TestInvalidConfigMissingTtl(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
domain: bacontest42.com
records:
    - host: bacontest42.com
      type: ALIAS
      content: pixie.porkbun.com
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	_, err = ReadFile(configFile)
	if err == nil {
		t.Fatal("expected error when a record is missing ttl field", err)
	}
}

func TestInvalidConfigInvalidTtl(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
domain: bacontest42.com
records:
    - host: bacontest42.com
      type: ALIAS
      ttl: 300
      content: pixie.porkbun.com
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	_, err = ReadFile(configFile)
	if err == nil {
		t.Fatal("expected error when a record has an invalid ttl", err)
	}
}

func TestInvalidConfigMissingContent(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
domain: bacontest42.com
records:
    - host: bacontest42.com
      type: ALIAS
      ttl: 600
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	_, err = ReadFile(configFile)
	if err == nil {
		t.Fatal("expected error when a record is missing content field", err)
	}
}

func TestInvalidConfigInvalidPriority(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
domain: bacontest42.com
records:
    - host: bacontest42.com
      type: ALIAS
      ttl: 600
      content: pixie.porkbun.com
	  priority: 20
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	_, err = ReadFile(configFile)
	if err == nil {
		t.Fatal("priority is not allowed on ALIAS records", err)
	}
}

func TestInvalidConfigCnameSameHost(t *testing.T) {
	configFile, err := SeedConfigToTempFile(`
domain: bacontest42.com
records:
    - host: '*.bacontest42.com'
      type: CNAME
      ttl: 600
      content: pixie.porkbun.com
    - type: MX
      host: "*.bacontest42.com"
      content: in1-smtp.messagingengine.com
      ttl: 600
      priority: 10
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	_, err = ReadFile(configFile)
	if err == nil {
		t.Fatal("cannot have a CNAME and another record for the same host", err)
	}
}

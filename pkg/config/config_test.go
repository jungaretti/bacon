package config

import (
	"testing"
)

func TestDeployMissingHost(t *testing.T) {
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

func TestDeployMissingType(t *testing.T) {
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

func TestDeployInvalidType(t *testing.T) {
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

func TestDeployMissingTtl(t *testing.T) {
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

func TestDeployInvalidTtl(t *testing.T) {
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

func TestDeployMissingContent(t *testing.T) {
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

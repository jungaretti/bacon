package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/providers/mock"
	"testing"
)

func TestDeploy(t *testing.T) {
	configFile, err := config.SeedConfigToTempFile(`
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
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	mockProvider := mock.NewMockProvider()

	err = deploy(mockProvider, configFile, true, true)
	if err != nil {
		t.Fatal("did not deploy records", err)
	}

	records, err := mockProvider.AllRecords("bacontest42.com")
	if err != nil {
		t.Fatal("could not fetch records after deployment", err)
	}

	if len(records) != 2 {
		t.Fatal("expected 2 records after deployment, got", len(records))
	}
}

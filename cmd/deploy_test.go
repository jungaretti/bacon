package cmd

import (
	console "bacon/pkg/providers/mock"
	"os"
	"testing"
)

func TestDeploy(t *testing.T) {
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
`)
	if err != nil {
		t.Fatal("could not seed config to temp file", err)
	}

	mockProvider := console.NewMockProvider()

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

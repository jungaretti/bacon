package cmd

import (
	"bacon/pkg/providers/console"
	"os"
	"testing"
)

func TestDeploy(t *testing.T) {
	configFile, err := seedConfigToTempFile(`
domain: plantbasedbacon.xyz
records:
    - host: plantbasedbacon.xyz
      type: ALIAS
      ttl: 600
      content: pixie.porkbun.com
    - host: '*.plantbasedbacon.xyz'
      type: CNAME
      ttl: 600
      content: pixie.porkbun.com
`)

	provider := console.NewConsoleProvider()
	mockApp := &App{
		Provider: provider,
	}

	err = deploy(mockApp, configFile, true, true)
	if err != nil {
		t.Fatal("did not deploy records", err)
	}

	recordsAfterDeployment, err := provider.AllRecords("plantbasedbacon.xyz")
	if err != nil {
		t.Fatal("could not fetch records after deployment", err)
	}

	if len(recordsAfterDeployment) != 2 {
		t.Fatal("expected 2 records after deployment, got", len(recordsAfterDeployment))
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

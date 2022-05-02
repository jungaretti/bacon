package client_test

import (
	"bacon/pkg/client"
	"bufio"
	"fmt"
	"os"
	"testing"
)

const (
	DOMAIN         = "borkbork.buzz"
	RECORD_TYPE    = "A"
	RECORD_HOST    = "borkbork.buzz"
	RECORD_CONTENT = "123.456.789.112"
	RECORD_TTL     = "600"
)

func TestReadConfig(t *testing.T) {
	temp, err := os.CreateTemp("", "baconconfig-")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	writer := bufio.NewWriter(temp)
	fmt.Fprintf(writer, "domain: %s\n", DOMAIN)
	fmt.Fprintf(writer, "records:\n")
	fmt.Fprintf(writer, "- type: %s\n", RECORD_TYPE)
	fmt.Fprintf(writer, "  host: %s\n", RECORD_HOST)
	fmt.Fprintf(writer, "  content: %s\n", RECORD_CONTENT)
	fmt.Fprintf(writer, "  ttl: \"%s\"\n", RECORD_TTL)
	err = temp.Close()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	config := client.Config{}
	err = client.ReadConfig(temp.Name(), &config)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestWriteConfig(t *testing.T) {
	config := client.Config{
		Domain: DOMAIN,
		Records: []client.Record{
			{
				Type:    RECORD_TYPE,
				Host:    RECORD_HOST,
				Content: RECORD_CONTENT,
				TTL:     RECORD_TTL,
			},
		},
	}

	temp, err := os.CreateTemp("", "baconconfig-")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = client.WriteConfig(temp.Name(), &config)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

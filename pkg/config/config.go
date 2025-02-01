package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Domain  string   `yaml:"domain"`
	Records []Record `yaml:"records"`
}

func (c Config) Validate() error {
	if c.Domain == "" {
		return fmt.Errorf("domain is required")
	}

	for _, record := range c.Records {
		err := record.Validate()
		if err != nil {
			return fmt.Errorf("record is invalid %v: %v", record, err)
		}
	}

	cnameHosts := make(map[string]bool)
	for _, record := range c.Records {
		if record.GetType() == "CNAME" {
			if _, ok := cnameHosts[record.GetName()]; ok {
				return fmt.Errorf("multiple CNAME records for host %s", record.GetName())
			}
			cnameHosts[record.GetName()] = true
		}
	}
	for _, record := range c.Records {
		if record.GetType() != "CNAME" && cnameHosts[record.GetName()] {
			return fmt.Errorf("non-CNAME record %v shares host with a CNAME record", record)
		}
	}

	return nil
}

func ReadFile(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(raw, &config)
	if err != nil {
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SeedConfigToTempFile(mockConfig string) (string, error) {
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

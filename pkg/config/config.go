package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Domain  string   `yaml:"domain"`
	Records []Record `yaml:"records"`
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

	err = ValidateConfiguration(config)
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

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Domain  string   `yaml:"domain"`
	Records []Record `yaml:"records"`
}

func ReadFile(configFile string) (*Config, error) {
	raw, err := os.ReadFile(configFile)
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

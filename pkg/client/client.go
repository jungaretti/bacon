package client

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Record struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Content  string `yaml:"content"`
	TTL      string `yaml:"ttl"`
	Priority string `yaml:"priority"`
}

type Config struct {
	Domain  string   `yaml:"domain"`
	Records []Record `yaml:"records"`
}

type Client interface {
	Name() string
	Ping() error
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(raw, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteConfig(filename string, config *Config) error {
	raw, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, raw, 0664)
}

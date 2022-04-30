package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

func ReadConfig(filename string, config *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(raw, &config)
	if err != nil {
		return err
	}

	return nil
}

func WriteConfig(filename string, config *Config) error {
	raw, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, raw, 0664)
}

func PostJson(url string, body interface{}) (*http.Response, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewReader(json))
}

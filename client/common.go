package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Ack struct {
	Ok      bool
	Message string
}

type Record struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Content  string `yaml:"content"`
	TTL      int    `yaml:"ttl"`
	Priority int    `yaml:"priority"`
}

type Config struct {
	Domain  string   `yaml:"domain"`
	Records []Record `yaml:"records"`
}

type Client interface {
	Ping() (*Ack, error)
	GetRecords(string) ([]Record, error)
	CreateRecord(string, *Record) error
	DeleteRecord(string, *Record) error
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

func postAndRead(url string, body interface{}) ([]byte, error) {
	enc, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(enc))
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}

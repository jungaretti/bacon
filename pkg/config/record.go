package config

import (
	"bacon/pkg/dns"
	"fmt"
)

type Record struct {
	Name string `yaml:"host"`
	Type string `yaml:"type"`
	Ttl  *int   `yaml:"ttl"`
	Data string `yaml:"content"`
}

func (r Record) GetName() string {
	return r.Name
}

func (r Record) GetType() string {
	return r.Type
}

func (r Record) GetTtl() string {
	if r.Ttl == nil {
		return ""
	}

	return fmt.Sprint(*r.Ttl)
}

func (r Record) GetData() string {
	return r.Data
}

var _ dns.Record = Record{}

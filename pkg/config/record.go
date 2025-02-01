package config

import (
	"bacon/pkg/dns"
	"fmt"
)

type Record struct {
	Name     string `yaml:"host"`
	Type     string `yaml:"type"`
	Ttl      int    `yaml:"ttl"`
	Data     string `yaml:"content"`
	Priority int    `yaml:"priority"`
}

func (r Record) GetName() string {
	return r.Name
}

func (r Record) GetType() string {
	return r.Type
}

func (r Record) GetTtl() string {
	return fmt.Sprint(r.Ttl)
}

func (r Record) GetData() string {
	return r.Data
}

func (r Record) GetPriority() string {
	if r.Priority == 0 {
		return ""
	}

	return fmt.Sprint(r.Priority)
}

var _ dns.Record = Record{}

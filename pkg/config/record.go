package config

import (
	"bacon/pkg/dns"
	"fmt"
	"strconv"
)

type ConfigRecord struct {
	Name string `yaml:"host"`
	Type string `yaml:"type"`
	Ttl  *int   `yaml:"ttl"`
	Data string `yaml:"content"`
}

func (r ConfigRecord) GetName() string {
	return r.Name
}

func (r ConfigRecord) GetType() string {
	return r.Type
}

func (r ConfigRecord) GetTtl() string {
	if r.Ttl == nil {
		return ""
	}

	return fmt.Sprint(*r.Ttl)
}

func (r ConfigRecord) GetData() string {
	return r.Data
}

var _ dns.Record = ConfigRecord{}

func ConfigFromRecord(r dns.Record) ConfigRecord {
	ttl, _ := strconv.Atoi(r.GetTtl())

	return ConfigRecord{
		Name: r.GetName(),
		Type: r.GetType(),
		Ttl:  &ttl,
		Data: r.GetData(),
	}
}

package config

import (
	"bacon/pkg/dns"
	"fmt"
)

type Record struct {
	Name string `yaml:"host"`
	Type string `yaml:"type"`
	Ttl  int    `yaml:"ttl"`
	Data string `yaml:"content"`
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

func (r Record) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("host is required")
	}

	if r.Type == "" {
		return fmt.Errorf("type is required")
	}

	allowedTypes := map[string]bool{
		"A":     true,
		"MX":    true,
		"CNAME": true,
		"ALIAS": true,
		"TXT":   true,
		"NS":    true,
		"AAAA":  true,
		"SRV":   true,
		"TLSA":  true,
		"CAA":   true,
		"HTTPS": true,
		"SVCB":  true,
	}
	if !allowedTypes[r.Type] {
		return fmt.Errorf("type must be one of A, MX, CNAME, ALIAS, TXT, NS, AAAA, SRV, TLSA, CAA, HTTPS, SVCB")
	}

	if r.Ttl < 600 {
		return fmt.Errorf("ttl must be at least 600")
	}

	if r.Data == "" {
		return fmt.Errorf("content is required")
	}

	return nil
}

var _ dns.Record = Record{}

package config

import (
	"bacon/pkg/dns"
	"fmt"
	"strings"
)

type Record struct {
	Name string `yaml:"host"`
	Type string `yaml:"type"`
	Ttl  int    `yaml:"ttl"`
	Data string `yaml:"content"`
}

const (
	TYPE_ALLOWLIST = "A, MX, CNAME, ALIAS, TXT, NS, AAAA, SRV, TLSA, CAA, HTTPS, SVCB"
)

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
	if !isTypeAllowed(r) {
		return fmt.Errorf("type must be one of %v", TYPE_ALLOWLIST)
	}

	if r.Ttl < 600 {
		return fmt.Errorf("ttl must be at least 600")
	}

	if r.Data == "" {
		return fmt.Errorf("content is required")
	}

	return nil
}

func isTypeAllowed(r Record) bool {
	allowedTypes := make(map[string]bool)
	for _, t := range strings.Split(TYPE_ALLOWLIST, ", ") {
		allowedTypes[t] = true
	}

	if _, ok := allowedTypes[r.Type]; !ok {
		return false
	}

	return true
}

var _ dns.Record = Record{}

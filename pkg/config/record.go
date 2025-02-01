package config

import (
	"bacon/pkg/dns"
	"fmt"
	"strings"
)

const (
	// Record types that are allowed. Found in Porkbun's API documentation.
	TYPE_ALLOWLIST = "A, MX, CNAME, ALIAS, TXT, NS, AAAA, SRV, TLSA, CAA, HTTPS, SVCB"
	// Record types that are allowed to have a priority. Found in Porkbun's web app.
	PRIORITY_ALLOWLIST = "MX, SRV"
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
	return fmt.Sprint(r.Priority)
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

	if r.Priority != 0 && !isPriorityAllowed(r) {
		return fmt.Errorf("priority must be one of %v", PRIORITY_ALLOWLIST)
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

func isPriorityAllowed(r Record) bool {
	allowedPriorities := make(map[string]bool)
	for _, t := range strings.Split(PRIORITY_ALLOWLIST, ", ") {
		allowedPriorities[t] = true
	}

	if _, ok := allowedPriorities[r.Type]; !ok {
		return false
	}

	return true
}

var _ dns.Record = Record{}

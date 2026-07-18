package config

import (
	"fmt"
	"slices"
	"strings"
)

var (
	// Record types that are allowed. Found in Porkbun's API documentation.
	allowedTypes = []string{"A", "MX", "CNAME", "ALIAS", "TXT", "AAAA", "SRV", "TLSA", "CAA", "HTTPS", "SVCB"}
	// Record types that are allowed to have a priority. Found in Porkbun's web app.
	allowedPriorityTypes = []string{"MX", "SRV"}
)

func ValidateConfiguration(config Config) error {
	if config.Domain == "" {
		return fmt.Errorf("domain is required")
	}

	for _, record := range config.Records {
		if err := ValidateRecord(record); err != nil {
			return fmt.Errorf("%v is invalid: %v", record, err)
		}
	}

	if err := configHasUniqueCnameHosts(config.Records); err != nil {
		return err
	}

	return nil
}

func ValidateRecord(record Record) error {
	if err := recordHasRequiredFields(record); err != nil {
		return err
	}

	if err := recordHasValidName(record); err != nil {
		return err
	}

	if err := recordHasValidType(record); err != nil {
		return err
	}

	if err := recordHasValidTtl(record); err != nil {
		return err
	}

	if err := recordHasValidPriority(record); err != nil {
		return err
	}

	return nil
}

func recordHasRequiredFields(record Record) error {
	if record.Name == "" {
		return fmt.Errorf("host is required")
	}

	if record.Type == "" {
		return fmt.Errorf("type is required")
	}

	if record.Data == "" {
		return fmt.Errorf("content is required")
	}

	return nil
}

func recordHasValidName(record Record) error {
	if strings.HasPrefix(record.Name, "_acme-challenge") {
		return fmt.Errorf("host must not begin with _acme-challenge")
	}

	return nil
}

func recordHasValidType(record Record) error {
	if !slices.Contains(allowedTypes, record.Type) {
		return fmt.Errorf("type must be one of %v", allowedTypes)
	}

	return nil
}

func recordHasValidTtl(record Record) error {
	if record.Ttl < 600 {
		return fmt.Errorf("ttl must be at least 600")
	}

	return nil
}

func recordHasValidPriority(record Record) error {
	if record.Priority == 0 {
		return nil
	}

	if !slices.Contains(allowedPriorityTypes, record.Type) {
		return fmt.Errorf("type must be one of %v to have priority", allowedPriorityTypes)
	}

	return nil
}

func configHasUniqueCnameHosts(records []Record) error {
	cnameHosts := make(map[string]bool)
	for _, record := range records {
		if record.Type == "CNAME" {
			if _, ok := cnameHosts[record.Name]; ok {
				return fmt.Errorf("multiple CNAME records exist for host %s", record.Name)
			}
			cnameHosts[record.Name] = true
		}
	}

	for _, record := range records {
		if record.Type == "CNAME" {
			continue
		}
		if cnameHosts[record.Name] {
			return fmt.Errorf("non-CNAME record %v shares host with a CNAME record", record)
		}
	}

	return nil
}

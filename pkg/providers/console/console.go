package console

import (
	"bacon/pkg/dns"
	"fmt"
)

type ConsoleProvider struct {
	records map[string]dns.Record
}

func (c ConsoleProvider) AllRecords(domain string) ([]dns.Record, error) {
	var result []dns.Record
	for _, record := range c.records {
		result = append(result, record)
	}
	return result, nil
}

func (c ConsoleProvider) CheckAuth() error {
	return nil
}

func (c ConsoleProvider) CreateRecord(domain string, record dns.Record) error {
	c.records[recordId(record)] = record
	return nil
}

func (c ConsoleProvider) DeleteRecord(domain string, record dns.Record) error {
	if _, ok := c.records[recordId(record)]; !ok {
		return fmt.Errorf("record not found")
	}

	delete(c.records, recordId(record))
	return nil
}

var _ dns.Provider = ConsoleProvider{}

func NewConsoleProvider() dns.Provider {
	return &ConsoleProvider{
		records: make(map[string]dns.Record),
	}
}

func recordId(record dns.Record) string {
	return fmt.Sprintf("%s-%s", record.GetName(), record.GetType())
}

package mock

import (
	"bacon/pkg/dns"
	"fmt"
)

type MockProvider struct {
	records map[string]dns.Record
}

func NewMockProvider() dns.Provider {
	return &MockProvider{
		records: make(map[string]dns.Record),
	}
}

func (c MockProvider) AllRecords(domain string) ([]dns.Record, error) {
	var result []dns.Record
	for _, record := range c.records {
		result = append(result, record)
	}
	return result, nil
}

func (c MockProvider) CheckAuth() error {
	return nil
}

func (c MockProvider) CreateRecord(domain string, record dns.Record) error {
	c.records[recordId(record)] = record
	return nil
}

func (c MockProvider) DeleteRecord(domain string, record dns.Record) error {
	if _, ok := c.records[recordId(record)]; !ok {
		return fmt.Errorf("record not found")
	}

	delete(c.records, recordId(record))
	return nil
}

func recordId(record dns.Record) string {
	return fmt.Sprintf("%s-%s", record.GetName(), record.GetType())
}

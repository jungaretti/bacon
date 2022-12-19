package mock

import (
	"bacon/pkg/dns"
	"fmt"
)

type MockProvider struct {
	Authenticated bool
	Records       *map[string]MockRecord
}

func (p MockProvider) AllRecords(domain string) ([]dns.Record, error) {
	records := make([]dns.Record, 0)
	for _, record := range *p.Records {
		records = append(records, record)
	}
	return records, nil
}

func (p MockProvider) CheckAuth() error {
	if p.Authenticated {
		return nil
	}

	return fmt.Errorf("unauthenticated")
}

func (p MockProvider) CreateRecord(_ string, record dns.Record) error {
	new := NewMockRecord(record)
	records := *p.Records
	records[new.GetId()] = new

	return nil
}

func (p MockProvider) DeleteRecord(domain string, record dns.Record) error {
	old := NewMockRecord(record)
	records := *p.Records
	delete(records, old.GetId())

	return nil
}

func NewMockProvider(authenticated bool) MockProvider {
	new := make(map[string]MockRecord, 0)

	return MockProvider{
		Authenticated: authenticated,
		Records:       &new,
	}
}

var _ dns.Provider = MockProvider{}

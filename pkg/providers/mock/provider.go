package mock

import "bacon/pkg/dns"

type MockProvider struct {
	Authenticated bool
	Records       map[MockRecord]bool
}

func (MockProvider) AllRecords(domain string) ([]dns.Record, error) {
	panic("unimplemented")
}

func (MockProvider) CheckAuth() error {
	panic("unimplemented")
}

func (MockProvider) CreateRecord(domain string, record dns.Record) error {
	panic("unimplemented")
}

func (MockProvider) DeleteRecord(domain string, record dns.Record) error {
	panic("unimplemented")
}

var _ dns.Provider = MockProvider{}

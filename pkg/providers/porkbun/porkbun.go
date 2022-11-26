package porkbun

import (
	"bacon/pkg/dns"
	"bacon/pkg/providers/porkbun/api"
	"bacon/pkg/providers/porkbun/record"
)

type PorkProvider struct {
	Api api.Api
}

func (p PorkProvider) CheckAuth() error {
	return p.Api.Ping()
}

func (p PorkProvider) AllRecords() ([]dns.Record, error) {
	porkRecords, err := p.Api.RetrieveRecords("")
	if err != nil {
		return nil, err
	}

	records := make([]dns.Record, len(porkRecords))
	for i, r := range porkRecords {
		records[i] = r
	}

	return records, nil
}

func (p PorkProvider) CreateRecord(newRecord dns.Record) error {
	porkRecord := record.Record{
		Type:    newRecord.GetType(),
		Name:    newRecord.GetHost(),
		Content: newRecord.GetContent(),
		TTL:     newRecord.GetTTL(),
	}

	_, err := p.Api.CreateRecord("", porkRecord)
	return err
}

func (p PorkProvider) DeleteRecord(dns.Record) error {
	// TODO: Find Porkbun ID of generic record and pass that here
	return p.Api.DeleteRecord("", "")
}

var _ dns.Provider = PorkProvider{}

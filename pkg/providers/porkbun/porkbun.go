package porkbun

import (
	"bacon/pkg/dns"
	"bacon/pkg/providers/porkbun/api"
	"bacon/pkg/providers/porkbun/record"
	"fmt"
)

type PorkProvider struct {
	Api api.Api
}

func (p PorkProvider) CheckAuth() error {
	return p.Api.Ping()
}

func (p PorkProvider) AllRecords(domain string) ([]dns.Record, error) {
	porkRecords, err := p.Api.RetrieveRecords(domain)
	if err != nil {
		return nil, err
	}

	records := make([]dns.Record, len(porkRecords))
	for i, r := range porkRecords {
		records[i] = r
	}

	return records, nil
}

func (p PorkProvider) CreateRecord(domain string, newRecord dns.Record) error {
	porkRecord := record.Record{
		Type:    newRecord.GetType(),
		Name:    newRecord.GetHost(),
		Content: newRecord.GetContent(),
		TTL:     newRecord.GetTTL(),
	}

	_, err := p.Api.CreateRecord(domain, porkRecord)
	return err
}

func (p PorkProvider) DeleteRecord(domain string, toDelete dns.Record) error {
	allRecords, err := p.Api.RetrieveRecords(domain)
	if err != nil {
		return err
	}

	var target *record.Record = nil
	for _, record := range allRecords {
		if dns.RecordEquals(toDelete, record) {
			target = &record
			break
		}
	}
	if target == nil {
		return fmt.Errorf("could not find a DNS record to delete")
	}

	return p.Api.DeleteRecord(domain, target.Id)
}

var _ dns.Provider = PorkProvider{}

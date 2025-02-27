package porkbun

import (
	"bacon/pkg/dns"
	"bacon/pkg/providers/porkbun/api"
	"fmt"
)

type PorkProvider struct {
	Api api.Api
}

func NewPorkbunProvider(apiKey, secretApiKey string) dns.Provider {
	return &PorkProvider{
		Api: api.Api{
			Auth: api.Auth{
				ApiKey:       apiKey,
				SecretApiKey: secretApiKey,
			},
			Throttler: *api.NewThrottler(1),
		},
	}
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
	porkRecord := api.Record{
		Name:    newRecord.GetName(),
		Type:    newRecord.GetType(),
		TTL:     newRecord.GetTtl(),
		Content: newRecord.GetData(),
	}

	if newRecord.GetPriority() != "0" {
		porkRecord.Priority = newRecord.GetPriority()
	}

	_, err := p.Api.CreateRecord(domain, porkRecord)
	return err
}

func (p PorkProvider) DeleteRecord(domain string, toDelete dns.Record) error {
	allRecords, err := p.Api.RetrieveRecords(domain)
	if err != nil {
		return err
	}

	var target *api.Record = nil
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

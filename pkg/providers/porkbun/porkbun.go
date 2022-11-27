package porkbun

import (
	"bacon/pkg/dns"
	"bacon/pkg/providers/porkbun/api"
	"bacon/pkg/providers/porkbun/record"
	"bacon/pkg/secrets"
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
		Name:    newRecord.GetName(),
		Type:    newRecord.GetType(),
		TTL:     newRecord.GetTtl(),
		Content: newRecord.GetData(),
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

func NewPorkbunProvider(secretsProvider secrets.Provider) dns.Provider {
	return &PorkProvider{
		Api: api.Api{
			Auth: api.Auth{
				ApiKey:       secretsProvider.Read(secrets.PorkbunApiKey),
				SecretApiKey: secretsProvider.Read(secrets.PorkbunSecretKey),
			},
		},
	}
}

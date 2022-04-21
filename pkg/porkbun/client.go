package porkbun

import (
	"bacon/pkg/client"
)

type PorkClient struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

type PorkRecord struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
	Notes    string `json:"notes"`
}

func (pork *PorkClient) GetName() string {
	return "Porkbun"
}

func (pork *PorkClient) Ping() (*client.Ack, error) {
	return pork.ping()
}

func (pork *PorkClient) GetRecords(domain string) ([]client.Record, error) {
	return pork.getRecords(domain)
}

func (pork *PorkClient) CreateRecord(domain string, record *client.Record) (*client.Ack, error) {
	return pork.createRecord(domain, record)
}

func (pork *PorkClient) DeleteRecord(domain string, id string) (*client.Ack, error) {
	return pork.deleteRecord(domain, id)
}

func (pork *PorkClient) SyncRecords(domain string, new []client.Record, create, delete bool) (*client.Ack, error) {
	return pork.syncRecords(domain, new, create, delete)
}

func ToRecord(porkRecord *PorkRecord) client.Record {
	return client.Record{
		Type:     porkRecord.Type,
		Host:     porkRecord.Name,
		Content:  porkRecord.Content,
		TTL:      porkRecord.TTL,
		Priority: porkRecord.Priority,
	}
}

func ToPorkRecord(record *client.Record) PorkRecord {
	return PorkRecord{
		Type:     record.Type,
		Name:     record.Host,
		Content:  record.Content,
		TTL:      record.TTL,
		Priority: record.Priority,
	}
}

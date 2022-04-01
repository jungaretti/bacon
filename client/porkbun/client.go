package porkbun

import (
	"bacon/client"
)

type PorkClient struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
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

package client

import (
	"fmt"
	"strconv"
)

const (
	PORK_PING          = "https://porkbun.com/api/json/v3/ping"
	PORK_GET_RECORDS   = "https://porkbun.com/api/json/v3/dns/retrieve/"
	PORK_CREATE_RECORD = "https://porkbun.com/api/json/v3/dns/create/"
	PORK_DELETE_RECORD = "https://porkbun.com/api/json/v3/dns/delete/"
)

type Pork struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

type porkRecord struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Host     string `json:"name"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
}

type porkCreation struct {
	Pork
	porkRecord
}

func (pork *Pork) Ping() (string, error) {
	raw, err := postAndRead(PORK_PING, pork)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func (pork *Pork) GetRecords(domain string) (string, error) {
	raw, err := postAndRead(PORK_GET_RECORDS+domain, pork)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func (pork *Pork) CreateRecord(domain string, record *Record) (string, error) {
	body := porkCreation{
		Pork:       *pork,
		porkRecord: *record.toPorkRecord(),
	}

	raw, err := postAndRead(PORK_CREATE_RECORD+domain, body)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func (pork *Pork) DeleteRecord(domain string, id string) (string, error) {
	raw, err := postAndRead(PORK_DELETE_RECORD+domain+"/"+id, pork)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func (record *Record) toPorkRecord() *porkRecord {
	return &porkRecord{
		Type:     record.Type,
		Host:     record.Host,
		Content:  record.Content,
		TTL:      fmt.Sprint(record.TTL),
		Priority: fmt.Sprint(record.Priority),
	}
}

func (porkRecord *porkRecord) toRecord() (*Record, error) {
	ttlInt, err := strconv.Atoi(porkRecord.TTL)
	if err != nil {
		return nil, err
	}
	priorityInt, err := strconv.Atoi(porkRecord.Priority)
	if err != nil {
		return nil, err
	}

	return &Record{
		Type:     porkRecord.Type,
		Host:     porkRecord.Host,
		Content:  porkRecord.Content,
		TTL:      ttlInt,
		Priority: priorityInt,
	}, nil
}

package client

import (
	"encoding/json"
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
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
	Notes    string `json:"notes"`
}

func (pork *Pork) Ping() (*Ack, error) {
	raw, err := postAndRead(PORK_PING, pork)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Status  string `json:"status"`
		YourIp  string `json:"yourIp"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, err
	}

	var ok bool
	switch resp.Status {
	case "SUCCESS":
		ok = true
	default:
		ok = false
	}

	var msg string
	switch ok {
	case true:
		msg = resp.YourIp
	default:
		msg = resp.Message
	}

	return &Ack{
		Ok:      ok,
		Message: msg,
	}, nil
}

// func GetRecords(string) ([]Record, error)        { return nil, nil }
// func CreateRecord(string, *Record) (*Ack, error) { return nil, nil }
// func DeleteRecord(string, *Record) (*Ack, error) { return nil, nil }

type porkCreation struct {
	Pork
	porkRecord
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
		Name:     record.Host,
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
		Host:     porkRecord.Name,
		Content:  porkRecord.Content,
		TTL:      ttlInt,
		Priority: priorityInt,
	}, nil
}

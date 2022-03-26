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

type porkBaseResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func parseStatus(status string) bool {
	switch status {
	case "SUCCESS":
		return true
	default:
		return false
	}
}

func parseMessage(status string, success string, failure string) string {
	switch parseStatus(status) {
	case true:
		return success
	default:
		return failure
	}
}

func (pork *Pork) Ping() (*Ack, error) {
	raw, err := postAndRead(PORK_PING, pork)
	if err != nil {
		return nil, err
	}

	var resp struct {
		YourIp string `json:"yourIp"`
		porkBaseResp
	}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, err
	}

	return &Ack{
		Ok:      parseStatus(resp.Status),
		Message: parseMessage(resp.Status, resp.YourIp, resp.Message),
	}, nil
}

func (pork *Pork) GetRecords(domain string) ([]Record, error) {
	raw, err := postAndRead(PORK_GET_RECORDS+domain, pork)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Records []porkRecord `json:"records"`
		porkBaseResp
	}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, err
	}
	if !parseStatus(resp.Status) {
		return nil, fmt.Errorf(resp.Message)
	}

	var records []Record
	for _, porkRec := range resp.Records {
		rec, err := porkRec.toRecord()
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

// func CreateRecord(string, *Record) (*Ack, error) { return nil, nil }
// func DeleteRecord(string, *Record) (*Ack, error) { return nil, nil }

type porkCreation struct {
	Pork
	porkRecord
}

func (pork *Pork) CreateRecord(domain string, record *Record) (string, error) {
	body := porkCreation{
		Pork:       *pork,
		porkRecord: record.toPorkRecord(),
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

func (record *Record) toPorkRecord() porkRecord {
	return porkRecord{
		Type:     record.Type,
		Name:     record.Host,
		Content:  record.Content,
		TTL:      fmt.Sprint(record.TTL),
		Priority: fmt.Sprint(record.Priority),
	}
}

func (porkRecord *porkRecord) toRecord() (Record, error) {
	ttlInt, err := strconv.Atoi(porkRecord.TTL)
	if err != nil {
		return Record{}, err
	}
	priorityInt, err := strconv.Atoi(porkRecord.Priority)
	if err != nil {
		return Record{}, err
	}

	return Record{
		Type:     porkRecord.Type,
		Host:     porkRecord.Name,
		Content:  porkRecord.Content,
		TTL:      ttlInt,
		Priority: priorityInt,
	}, nil
}

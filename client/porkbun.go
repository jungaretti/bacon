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

func (pork *Pork) getPorkRecords(domain string) ([]porkRecord, error) {
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

	return resp.Records, nil
}

func (pork *Pork) GetRecords(domain string) ([]Record, error) {
	porkRecords, err := pork.getPorkRecords(domain)
	if err != nil {
		return nil, err
	}

	var records []Record
	for _, porkRec := range porkRecords {
		rec, err := porkRec.toRecord()
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (pork *Pork) CreateRecord(domain string, record *Record) (*Ack, error) {
	create := struct {
		Pork
		porkRecord
	}{
		Pork:       *pork,
		porkRecord: record.toPorkRecord(),
	}

	raw, err := postAndRead(PORK_CREATE_RECORD+domain, create)
	if err != nil {
		return nil, err
	}

	var resp struct {
		// TODO: Report this bug, should be a string
		Id int `json:"id"`
		porkBaseResp
	}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, err
	}
	if !parseStatus(resp.Status) {
		return nil, fmt.Errorf(resp.Message)
	}

	return &Ack{
		Ok:      parseStatus(resp.Status),
		Message: parseMessage(resp.Status, fmt.Sprint(resp.Id), resp.Message),
	}, nil
}

func (pork *Pork) DeleteRecord(domain string, id string) (*Ack, error) {
	raw, err := postAndRead(PORK_DELETE_RECORD+domain+"/"+id, pork)
	if err != nil {
		return nil, err
	}

	resp := porkBaseResp{}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, err
	}
	if !parseStatus(resp.Status) {
		return nil, fmt.Errorf(resp.Message)
	}

	return &Ack{
		Ok:      parseStatus(resp.Status),
		Message: parseMessage(resp.Status, fmt.Sprint(id), resp.Message),
	}, nil
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

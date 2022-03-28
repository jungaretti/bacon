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

type PorkClient struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

func (pork *PorkClient) Ping() (*Ack, error) {
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

func (pork *PorkClient) GetRecords(domain string) ([]Record, error) {
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

func (pork *PorkClient) CreateRecord(domain string, record *Record) (*Ack, error) {
	create := struct {
		PorkClient
		porkRecord
	}{
		PorkClient: *pork,
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

func (pork *PorkClient) DeleteRecord(domain string, id string) (*Ack, error) {
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

func (pork *PorkClient) SyncRecords(domain string, new []Record, create, delete bool) (*Ack, error) {
	old, err := pork.GetRecords(domain)
	if err != nil {
		return nil, err
	}

	toDelete := difference(old, new)
	toCreate := difference(new, old)

	porkRecords, err := pork.getPorkRecords(domain)
	if err != nil {
		return nil, err
	}

	// Delete records
	for _, record := range toDelete {
		id, err := findId(&record, porkRecords)
		if err != nil {
			return nil, err
		}

		if delete {
			ack, err := pork.DeleteRecord(domain, id)
			if err != nil {
				return nil, err
			}
			if !ack.Ok {
				return nil, fmt.Errorf(ack.Message)
			}
		} else {
			fmt.Printf("Would delete record %s: ", id)
			fmt.Println(record)
		}
	}

	// Create records
	for _, record := range toCreate {
		if create {
			ack, err := pork.CreateRecord(domain, &record)
			if err != nil {
				return nil, err
			}
			if !ack.Ok {
				return nil, fmt.Errorf(ack.Message)
			}
		} else {
			fmt.Printf("Would create record")
			fmt.Println(record)
		}
	}

	return &Ack{
		Ok:      true,
		Message: "Synced records with Porkbun",
	}, nil
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

func (record *Record) toPorkRecord() porkRecord {
	return porkRecord{
		Type:     record.Type,
		Name:     record.Host,
		Content:  record.Content,
		TTL:      fmt.Sprint(record.TTL),
		Priority: fmt.Sprint(record.Priority),
	}
}

type porkBaseResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (pork *PorkClient) getPorkRecords(domain string) ([]porkRecord, error) {
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

func findId(target *Record, all []porkRecord) (string, error) {
	for _, porkRec := range all {
		rec, err := porkRec.toRecord()
		if err != nil {
			return "", err
		}

		if rec == *target {
			return porkRec.Id, nil
		}
	}

	return "", fmt.Errorf("didn't find a matching Porkbun record")
}

func difference(a, b []Record) (diff []Record) {
	m := make(map[Record]bool)

	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

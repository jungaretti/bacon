package porkbun

import (
	"bacon/pkg/client"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	PORK_PING          = "https://porkbun.com/api/json/v3/ping"
	PORK_GET_RECORDS   = "https://porkbun.com/api/json/v3/dns/retrieve/"
	PORK_CREATE_RECORD = "https://porkbun.com/api/json/v3/dns/create/"
	PORK_DELETE_RECORD = "https://porkbun.com/api/json/v3/dns/delete/"
)

type porkBaseResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (pork *PorkClient) ping() (*client.Ack, error) {
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

	return &client.Ack{
		Ok:      parseStatus(resp.Status),
		Message: parseMessage(resp.Status, resp.YourIp, resp.Message),
	}, nil
}

func (pork *PorkClient) getPorkRecords(domain string) ([]PorkRecord, error) {
	raw, err := postAndRead(PORK_GET_RECORDS+domain, pork)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Records []PorkRecord `json:"records"`
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

func (pork *PorkClient) getRecords(domain string) ([]client.Record, error) {
	porkRecords, err := pork.getPorkRecords(domain)
	if err != nil {
		return nil, err
	}

	var records []client.Record
	for _, porkRec := range porkRecords {
		rec := ToRecord(&porkRec)
		records = append(records, rec)
	}

	return records, nil
}

func (pork *PorkClient) createRecord(domain string, record *client.Record) (*client.Ack, error) {
	create := struct {
		PorkClient
		PorkRecord
	}{
		PorkClient: *pork,
		PorkRecord: ToPorkRecord(record),
	}

	create.Name = trimDomainFromHost(domain, create.Name)

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

	return &client.Ack{
		Ok:      parseStatus(resp.Status),
		Message: parseMessage(resp.Status, fmt.Sprint(resp.Id), resp.Message),
	}, nil
}

func (pork *PorkClient) deleteRecord(domain string, id string) (*client.Ack, error) {
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

	return &client.Ack{
		Ok:      parseStatus(resp.Status),
		Message: parseMessage(resp.Status, fmt.Sprint(id), resp.Message),
	}, nil
}

func (pork *PorkClient) syncRecords(domain string, new []client.Record, create, delete bool) (*client.Ack, error) {
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

	return &client.Ack{
		Ok:      true,
		Message: "Synced records with Porkbun",
	}, nil
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

func findId(target *client.Record, all []PorkRecord) (string, error) {
	for _, porkRec := range all {
		rec := ToRecord(&porkRec)
		if rec == *target {
			return porkRec.Id, nil
		}
	}

	return "", fmt.Errorf("didn't find a matching Porkbun record")
}

func trimDomainFromHost(domain, host string) string {
	if host == domain {
		return ""
	} else {
		return strings.Replace(host, "."+domain, "", 1)
	}
}

func difference(a, b []client.Record) (diff []client.Record) {
	m := make(map[client.Record]bool)

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

func postAndRead(url string, body interface{}) ([]byte, error) {
	enc, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(enc))
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}

package api

import (
	"fmt"
	"strings"
)

const (
	PING     string = "https://api.porkbun.com/api/json/v3/ping"
	RETRIEVE string = "https://api.porkbun.com/api/json/v3/dns/retrieve"
	CREATE   string = "https://api.porkbun.com/api/json/v3/dns/create"
	DELETE   string = "https://api.porkbun.com/api/json/v3/dns/delete"
)

type Api struct {
	Auth      Auth
	Throttler Throttler
}

func (p Api) Ping() error {
	type pingRes struct {
		baseRes
		YourIp string `json:"yourIp"`
	}

	response := pingRes{}
	p.Throttler.WaitForPermit()
	err := makeRequest(PING, p.Auth, &response)
	if err != nil {
		return err
	}

	return nil
}

func (p Api) RetrieveRecords(domain string) ([]Record, error) {
	type listRes struct {
		baseRes
		Records []Record `json:"records"`
	}

	response := listRes{}
	p.Throttler.WaitForPermit()
	err := makeRequest(RETRIEVE+"/"+domain, p.Auth, &response)
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0)
	for _, record := range response.Records {
		if record.isIgnored() {
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

func (p Api) CreateRecord(domain string, toCreate Record) (string, error) {
	type createReq struct {
		Auth
		Record
	}

	type createRes struct {
		baseRes
		Id int `json:"id"`
	}

	if toCreate.isIgnored() {
		return "", fmt.Errorf("cannot create an ignored record: %s", toCreate)
	}

	request := createReq{
		Auth:   p.Auth,
		Record: toCreate,
	}
	request.Name = trimDomain(toCreate.Name, domain)

	response := createRes{}
	p.Throttler.WaitForPermit()
	err := makeRequest(CREATE+"/"+domain, request, &response)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(response.Id), nil
}

func (p Api) DeleteRecord(domain string, id string) error {
	response := baseRes{}
	p.Throttler.WaitForPermit()
	err := makeRequest(DELETE+"/"+domain+"/"+id, p.Auth, &response)
	if err != nil {
		return err
	}

	return nil
}

// Trims a root domain from a longer subdomain. For example, trims
// host.example.com to host. If the subdomain is example.com, then
// returns an empty string
func trimDomain(name string, domain string) string {
	if name == domain {
		return ""
	}

	return strings.Replace(name, "."+domain, "", 1)
}

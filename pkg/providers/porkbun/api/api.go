package api

import (
	porkbun "bacon/pkg/providers/porkbun/record"
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
	Auth Auth
}

func (p Api) Ping() error {
	type pingRes struct {
		baseRes
		YourIp string `json:"yourIp"`
	}

	response := pingRes{}
	err := makeRequest(PING, p.Auth, &response)
	if err != nil {
		return err
	}

	return nil
}

func (p Api) RetrieveRecords(domain string) ([]porkbun.Record, error) {
	type listRes struct {
		baseRes
		Records []porkbun.Record `json:"records"`
	}

	response := listRes{}
	err := makeRequest(RETRIEVE+"/"+domain, p.Auth, &response)
	if err != nil {
		return nil, err
	}

	records := ignoreRecords(response.Records)
	return records, nil
}

func (p Api) CreateRecord(domain string, toCreate porkbun.Record) (string, error) {
	type createReq struct {
		Auth
		porkbun.Record
	}

	type createRes struct {
		baseRes
		Id int `json:"id"`
	}

	if isIgnored(toCreate) {
		return "", fmt.Errorf("cannot create an ignored record: %s", toCreate)
	}

	request := createReq{
		Auth:   p.Auth,
		Record: toCreate,
	}
	request.Name = trimDomain(toCreate.Name, domain)

	response := createRes{}
	err := makeRequest(CREATE+"/"+domain, request, &response)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(response.Id), nil
}

func (p Api) DeleteRecord(domain string, id string) error {
	response := baseRes{}
	err := makeRequest(DELETE+"/"+domain+"/"+id, p.Auth, &response)
	if err != nil {
		return err
	}

	return nil
}

func isIgnored(record porkbun.Record) bool {
	if record.Type == "NS" {
		return true
	}
	if strings.HasPrefix(record.Name, "_acme-challenge") {
		return true
	}

	return false
}

func ignoreRecords(input []porkbun.Record) []porkbun.Record {
	records := make([]porkbun.Record, 0)
	for _, record := range input {
		if isIgnored(record) {
			continue
		}

		records = append(records, record)
	}

	return records
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

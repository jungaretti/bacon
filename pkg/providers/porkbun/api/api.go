package api

import porkbun "bacon/pkg/providers/porkbun/record"

const (
	PING     string = "https://porkbun.com/api/json/v3/ping"
	RETRIEVE string = "https://porkbun.com/api/json/v3/dns/retrieve"
	CREATE   string = "https://porkbun.com/api/json/v3/dns/create"
	DELETE   string = "https://porkbun.com/api/json/v3/dns/delete"
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

	return response.Records, nil
}

func (p Api) CreateRecord(domain string, toCreate porkbun.Record) (string, error) {
	type createReq struct {
		Auth
		porkbun.Record
	}

	type createRes struct {
		baseRes
		Id string `json:"id"`
	}

	request := createReq{
		Auth:   p.Auth,
		Record: toCreate,
	}

	response := createRes{}
	err := makeRequest(CREATE+"/"+domain, request, &response)
	if err != nil {
		return "", err
	}

	return response.Id, nil
}

func (p Api) DeleteRecord(domain string, id string) error {
	response := baseRes{}
	err := makeRequest(DELETE+"/"+domain+"/"+id, p.Auth, &response)
	if err != nil {
		return err
	}

	return nil
}

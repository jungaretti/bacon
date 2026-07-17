package porkbun

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

type Client struct {
	Auth      Auth
	Throttler Throttler
}

func NewClient(apiKey, secretApiKey string) *Client {
	return &Client{
		Auth:      Auth{ApiKey: apiKey, SecretApiKey: secretApiKey},
		Throttler: *NewThrottler(1),
	}
}

func (c Client) Ping() error {
	type pingRes struct {
		baseRes
		YourIp string `json:"yourIp"`
	}

	response := pingRes{}
	c.Throttler.WaitForPermit()
	return makeRequest(PING, c.Auth, &response)
}

func (c Client) AllRecords(domain string) ([]Record, error) {
	type listRes struct {
		baseRes
		Records []Record `json:"records"`
	}

	response := listRes{}
	c.Throttler.WaitForPermit()
	err := makeRequest(RETRIEVE+"/"+domain, c.Auth, &response)
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

func (c Client) CreateRecord(domain string, record Record) error {
	type createReq struct {
		Auth
		Record
	}

	type createRes struct {
		baseRes
		Id int `json:"id"`
	}

	if record.isIgnored() {
		return fmt.Errorf("cannot create an ignored record: %s", record)
	}

	request := createReq{Auth: c.Auth, Record: record}
	request.Name = trimDomain(record.Name, domain)

	response := createRes{}
	c.Throttler.WaitForPermit()
	return makeRequest(CREATE+"/"+domain, request, &response)
}

func (c Client) DeleteRecord(domain string, record Record) error {
	c.Throttler.WaitForPermit()
	response := baseRes{}
	return makeRequest(DELETE+"/"+domain+"/"+record.Id, c.Auth, &response)
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

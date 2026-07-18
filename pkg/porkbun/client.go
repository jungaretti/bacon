package porkbun

import (
	"fmt"
	"strings"
	"time"
)

const (
	pingUrl     string = "https://api.porkbun.com/api/json/v3/ping"
	retrieveUrl string = "https://api.porkbun.com/api/json/v3/dns/retrieve"
	createUrl   string = "https://api.porkbun.com/api/json/v3/dns/create"
	deleteUrl   string = "https://api.porkbun.com/api/json/v3/dns/delete"
)

// Porkbun returns 5XX errors if we send requests too quickly
const rateLimit = time.Second

type Auth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

type Client struct {
	Auth   Auth
	ticker *time.Ticker
}

func NewClient(apiKey, secretApiKey string) *Client {
	return &Client{
		Auth:   Auth{ApiKey: apiKey, SecretApiKey: secretApiKey},
		ticker: time.NewTicker(rateLimit),
	}
}

func (c Client) Ping() error {
	type pingRes struct {
		baseRes
		YourIp string `json:"yourIp"`
	}

	response := pingRes{}
	return makeRequest(c.ticker, pingUrl, c.Auth, &response)
}

func (c Client) AllRecords(domain string) ([]Record, error) {
	type listRes struct {
		baseRes
		Records []Record `json:"records"`
	}

	response := listRes{}
	err := makeRequest(c.ticker, retrieveUrl+"/"+domain, c.Auth, &response)
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
	return makeRequest(c.ticker, createUrl+"/"+domain, request, &response)
}

func (c Client) DeleteRecord(domain string, record Record) error {
	response := baseRes{}
	return makeRequest(c.ticker, deleteUrl+"/"+domain+"/"+record.Id, c.Auth, &response)
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

package porkbun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

const baseUrl = "https://api.porkbun.com/api/json/v3"

const (
	pingUrl     string = baseUrl + "/ping"
	retrieveUrl string = baseUrl + "/dns/retrieve"
	createUrl   string = baseUrl + "/dns/create"
	deleteUrl   string = baseUrl + "/dns/delete"
)

// Porkbun returns 5XX errors if we send requests too quickly
const rateLimit = time.Second

type Auth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

type Client struct {
	Auth        Auth
	nextRequest time.Time
}

func NewClient(apiKey, secretApiKey string) *Client {
	return &Client{
		Auth: Auth{ApiKey: apiKey, SecretApiKey: secretApiKey},
	}
}

func (client *Client) Ping() error {
	type pingRes struct {
		baseRes
		YourIp string `json:"yourIp"`
	}

	response := pingRes{}
	return client.post(pingUrl, client.Auth, &response)
}

func (client *Client) AllRecords(domain string) ([]Record, error) {
	type listRes struct {
		baseRes
		Records []Record `json:"records"`
	}

	response := listRes{}
	err := client.post(retrieveUrl+"/"+domain, client.Auth, &response)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(response.Records, Record.isIgnored), nil
}

func (client *Client) CreateRecord(domain string, record Record) (string, error) {
	type createReq struct {
		Auth
		Record
	}

	type createRes struct {
		baseRes
		Id int `json:"id"`
	}

	if record.isIgnored() {
		return "", fmt.Errorf("cannot create an ignored record: %s", record)
	}

	request := createReq{Auth: client.Auth, Record: record}
	request.Name = trimDomain(record.Name, domain)

	response := createRes{}
	err := client.post(createUrl+"/"+domain, request, &response)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(response.Id), nil
}

func (client *Client) DeleteRecord(domain string, record Record) error {
	response := baseRes{}
	return client.post(deleteUrl+"/"+domain+"/"+record.Id, client.Auth, &response)
}

// Trims a root domain from a longer subdomain. For example, trims
// host.example.com to host. If the subdomain is example.com, then
// returns an empty string
func trimDomain(name string, domain string) string {
	if name == domain {
		return ""
	}

	return strings.TrimSuffix(name, "."+domain)
}

type checkable interface {
	checkStatus() error
}

type baseRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (r baseRes) checkStatus() error {
	if r.Status == "SUCCESS" {
		return nil
	}

	return fmt.Errorf("unsuccessful Porkbun request: %s", r.Message)
}

func (client *Client) throttle() {
	time.Sleep(time.Until(client.nextRequest))
	client.nextRequest = time.Now().Add(rateLimit)
}

func (client *Client) post(url string, request any, response checkable) error {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	client.throttle()

	res, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("making POST request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("received non-success status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return fmt.Errorf("unmarshaling JSON body: %v", err)
	}

	return response.checkStatus()
}

package porkbun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	editUrl     string = baseUrl + "/dns/edit"
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
	_, err := post[struct{}](client, pingUrl, client.Auth)
	return err
}

func (client *Client) AllRecords(domain string) ([]Record, error) {
	type listRes struct {
		Records []Record `json:"records"`
	}

	response, err := post[listRes](client, retrieveUrl+"/"+domain, client.Auth)
	if err != nil {
		return nil, fmt.Errorf("retrieve records failed: %v", err)
	}

	records := slices.DeleteFunc(response.Records, Record.isIgnored)
	for i, record := range records {
		records[i] = record.NormalizePriority()
	}

	return records, nil
}

func (client *Client) CreateRecord(domain string, record Record) (string, error) {
	type createReq struct {
		Auth
		Record
	}

	type createRes struct {
		Id int `json:"id"`
	}

	if record.isIgnored() {
		return "", fmt.Errorf("cannot create an ignored record: %s", record)
	}

	request := createReq{Auth: client.Auth, Record: record}
	request.Name = trimDomain(record.Name, domain)

	response, err := post[createRes](client, createUrl+"/"+domain, request)
	if err != nil {
		return "", fmt.Errorf("create record failed: %v", err)
	}

	return strconv.Itoa(response.Id), nil
}

func (client *Client) EditRecord(domain string, record Record) error {
	type editReq struct {
		Auth
		Record
	}

	if record.isIgnored() {
		return fmt.Errorf("cannot edit ignored record: %s", record)
	}

	request := editReq{Auth: client.Auth, Record: record}
	request.Name = trimDomain(record.Name, domain)

	_, err := post[struct{}](client, editUrl+"/"+domain+"/"+record.Id, request)
	if err != nil {
		return fmt.Errorf("update record failed: %v", err)
	}

	return nil
}

func (client *Client) DeleteRecord(domain string, record Record) error {
	_, err := post[struct{}](client, deleteUrl+"/"+domain+"/"+record.Id, client.Auth)
	if err != nil {
		return fmt.Errorf("delete record failed: %v", err)
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

	return strings.TrimSuffix(name, "."+domain)
}

func (client *Client) throttle() {
	time.Sleep(time.Until(client.nextRequest))
	client.nextRequest = time.Now().Add(rateLimit)
}

func post[T any](client *Client, url string, request any) (T, error) {
	type baseRes struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var response T

	jsonData, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	client.throttle()

	res, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return response, fmt.Errorf("making POST request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		return response, fmt.Errorf("validation error")
	} else if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, fmt.Errorf("reading response body: %v", err)
	}

	base := baseRes{}
	if err = json.Unmarshal(body, &base); err != nil {
		return response, fmt.Errorf("unmarshaling JSON body: %v", err)
	}
	if base.Status != "SUCCESS" {
		return response, fmt.Errorf("unsuccessful Porkbun request: %s", base.Message)
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("unmarshaling JSON body: %v", err)
	}

	return response, nil
}

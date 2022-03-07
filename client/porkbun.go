package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	PORK_PING          = "https://porkbun.com/api/json/v3/ping"
	PORK_GET_RECORDS   = "https://porkbun.com/api/json/v3/dns/retrieve/"
	PORK_CREATE_RECORD = "https://porkbun.com/api/json/v3/dns/create/"
	PORK_DELETE_RECORD = "https://porkbun.com/api/json/v3/dns/delete/"
)

type Auth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

type Record struct {
	Type     string
	Host     string
	Content  string
	TTL      string
	Priority string
	Notes    string
}

func PingJSON(auth *Auth) ([]byte, error) {
	raw, err := postAndRead(PORK_PING, auth)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func GetRecordsJSON(auth *Auth, domain string) ([]byte, error) {
	raw, err := postAndRead(PORK_GET_RECORDS+domain, auth)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func CreateRecordJSON(auth *Auth, domain string, record *Record) ([]byte, error) {
	body := struct {
		ApiKey       string `json:"apikey"`
		SecretApiKey string `json:"secretapikey"`
		Host         string `json:"name"`
		Type         string `json:"type"`
		Content      string `json:"content"`
		TTL          string `json:"ttl"`
		Priority     string `json:"prio"`
	}{
		ApiKey:       auth.ApiKey,
		SecretApiKey: auth.SecretApiKey,
		Host:         record.Host,
		Type:         record.Type,
		Content:      record.Content,
		TTL:          record.TTL,
		Priority:     record.Priority,
	}

	raw, err := postAndRead(PORK_CREATE_RECORD+domain, body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func DeleteRecordJSON(auth *Auth, domain string, id string) ([]byte, error) {
	raw, err := postAndRead(PORK_DELETE_RECORD+domain+"/"+id, auth)
	if err != nil {
		return nil, err
	}

	return raw, nil
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

package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var PK_PING = "https://porkbun.com/api/json/v3/ping"
var PK_DNS_RETRIEVE = "https://porkbun.com/api/json/v3/dns/retrieve/"

type Auth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

type Record struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Ttl      string `json:"ttl"`
	Priority string `json:"prio"`
	Notes    string `json:"notes"`
}

type PingRes struct {
	Status   string `json:"status"`
	ClientIp string `json:"yourIp"`
}

type RecordRes struct {
	Status  string    `json:"status"`
	Records *[]Record `json:"records"`
}

func Ping(auth Auth) (*PingRes, error) {
	body, err := postAndDecode(auth, PK_PING)
	if err != nil {
		return nil, err
	}

	res := PingRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func RetrieveRecords(auth Auth, domain string) (*[]Record, error) {
	body, err := postAndDecode(auth, PK_DNS_RETRIEVE+domain)
	if err != nil {
		return nil, err
	}

	res := RecordRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res.Records, nil
}

func postAndDecode(auth Auth, url string) ([]byte, error) {
	json, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

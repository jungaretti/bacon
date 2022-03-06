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

type Ack struct {
	Status   string `json:"status"`
	ClientIp string `json:"yourIp"`
}

func (a *Ack) Ok() bool {
	switch a.Status {
	case "SUCCESS":
		return true
	default:
		return false
	}
}

func Ping(auth Auth) (*Ack, error) {
	body, err := postAndDecode(auth, PK_PING)
	if err != nil {
		return nil, err
	}

	ack := Ack{}
	err = json.Unmarshal(body, &ack)
	if err != nil {
		return nil, err
	}

	return &ack, nil
}

func RetrieveRecords(auth Auth, domain string) (string, error) {
	body, err := postAndDecode(auth, PK_DNS_RETRIEVE+domain)
	if err != nil {
		return "", err
	}

	return string(body), nil
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

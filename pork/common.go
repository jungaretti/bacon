package pork

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

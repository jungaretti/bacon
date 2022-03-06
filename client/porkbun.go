package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var PORK_PING = "https://porkbun.com/api/json/v3/ping"
var PORK_GET_RECORDS = "https://porkbun.com/api/json/v3/dns/retrieve/"
var PORK_CREATE_RECORD = "https://porkbun.com/api/json/v3/dns/create/"
var PORK_DELETE_RECORD = "https://porkbun.com/api/json/v3/dns/delete/"

type Auth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type Ack struct {
	Success bool
	Message string
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

type pingRes struct {
	Status   string `json:"status"`
	ClientIp string `json:"yourIp"`
	Message  string `josn:"message"`
}

type recordsRes struct {
	Status  string   `json:"status"`
	Records []Record `json:"records"`
}

func Ping(auth *Auth) (*Ack, error) {
	res := pingRes{}
	err := postAndRead(auth, PORK_PING, &res)
	if err != nil {
		return nil, err
	}

	ack := Ack{
		Success: parseStatus(res.Status),
		Message: res.Message,
	}

	return &ack, nil
}

func GetRecords(auth *Auth, domain string) (*[]Record, error) {
	res := recordsRes{}
	err := postAndRead(auth, PORK_GET_RECORDS+domain, &res)
	if err != nil {
		return nil, err
	}

	return &res.Records, nil
}

func CreateRecord(auth *Auth, domain string, record *Record) (bool, error) {
	return false, fmt.Errorf("couldn't create record")
}

func DeleteRecord(auth *Auth, domain string, id string) (bool, error) {
	return false, fmt.Errorf("couldn't delete record")
}

func postAndRead(body interface{}, url string, value interface{}) error {
	encoded, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, &value)
}

func parseStatus(status string) bool {
	switch status {
	case "SUCCESS":
		return true
	default:
		return false
	}
}

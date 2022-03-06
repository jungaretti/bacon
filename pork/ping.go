package pork

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Auth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

func Ping() (string, error) {
	auth := Auth{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	data, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("https://porkbun.com/api/json/v3/ping", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

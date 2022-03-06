package pork

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Ping(auth Auth) (string, error) {
	json, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(PORKBUN_PING, "application/json", bytes.NewBuffer(json))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

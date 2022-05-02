package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func PostJson(url string, body interface{}) (*http.Response, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewReader(json))
}

func PostJsonAndRead(url string, body interface{}) ([]byte, error) {
	res, err := PostJson(url, body)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}

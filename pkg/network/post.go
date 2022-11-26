package network

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func PostJson(url string, body interface{}) (*http.Response, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewReader(json))
}

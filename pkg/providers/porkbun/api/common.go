package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type checkable interface {
	checkStatus() error
}

type baseRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (r baseRes) checkStatus() error {
	if r.Status == "SUCCESS" {
		return nil
	}

	return fmt.Errorf("unsuccessful Porkbun request: %s", r.Message)
}

var _ checkable = baseRes{}

func makeRequest(url string, req interface{}, out checkable) error {
	res, err := postJson(url, req)
	if err != nil {
		return fmt.Errorf("making POST request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("received non-success status code: %d. Body: %s", res.StatusCode, string(body))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading JSON body: %v", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("closing response: %v", err)
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return fmt.Errorf("unmarshaling JSON body: %v", err)
	}

	if err = out.checkStatus(); err != nil {
		return err
	}

	return nil
}

func postJson(url string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewReader(jsonData))
}

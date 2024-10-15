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
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("making POST request: %v", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("received non-success status code: %d", res.StatusCode)
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

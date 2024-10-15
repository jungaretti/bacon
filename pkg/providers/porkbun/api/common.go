package api

import (
	"bacon/pkg/network"
	"encoding/json"
	"fmt"
	"io"
	"time"
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

func makeRequestWithBackoff(url string, req interface{}, out checkable) error {
	baseDelay := time.Duration(2 * float64(time.Second))
	maxDelay := time.Duration(4 * float64(time.Second))
	return RetryWithBackoff(5, baseDelay, maxDelay, func() error { return makeRequest(url, req, out) })
}

func makeRequest(url string, req interface{}, out checkable) error {
	res, err := network.PostJson(url, req)
	if err != nil {
		return fmt.Errorf("making POST request: %v", err)
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

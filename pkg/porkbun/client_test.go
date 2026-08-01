package porkbun

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPostSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`{"status": "SUCCESS", "yourIp": "1.2.3.4"}`))
	}))
	defer server.Close()

	type pingRes struct {
		YourIp string `json:"yourIp"`
	}

	client := NewClient("key", "secret")
	response, err := post[pingRes](client, server.URL, client.Auth)
	if err != nil {
		t.Errorf("status mismatch: error when successful: %v", err)
	}
	if response.YourIp != "1.2.3.4" {
		t.Errorf("unexpected response: %+v", response)
	}
}

func TestPostRetryAfterError(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		if requests == 1 {
			writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		writer.Write([]byte(`{"status": "SUCCESS"}`))
	}))
	defer server.Close()

	client := NewClient("key", "secret")
	client.rateLimit = 10 * time.Millisecond

	_, err := post[struct{}](client, server.URL, client.Auth)
	if err != nil {
		t.Errorf("status mismatch: error when successful: %v", err)
	}
	if requests != 2 {
		t.Errorf("unexpected number of requests: %d", requests)
	}
}

func TestPostRetryRespectsLimit(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		writer.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient("key", "secret")
	client.rateLimit = 10 * time.Millisecond

	_, err := post[struct{}](client, server.URL, client.Auth)
	if err == nil {
		t.Error("status mismatch: no error when unsuccessful")
	}
	if requests != client.retryLimit {
		t.Errorf("unexpected number of requests: %d", requests)
	}
}

func TestPostFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`{"status": "FAILURE", "message": "random error message"}`))
	}))
	defer server.Close()

	client := NewClient("key", "secret")
	_, err := post[struct{}](client, server.URL, client.Auth)
	if err == nil {
		t.Error("status mismatch: no error when unsuccessful")
	}
}

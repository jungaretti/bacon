package porkbun

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

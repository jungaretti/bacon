package helpers

import (
	"bacon/pkg/helpers"
	"testing"
)

func TestPostJson(t *testing.T) {
	body := struct {
		Name string `json:"name"`
	}{
		Name: "Steve Jobs",
	}

	resp, err := helpers.PostJson("https://ptsv2.com/t/tg7o8-1651468588/post", body)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if resp.StatusCode != 200 {
		t.Log(resp.StatusCode)
		t.FailNow()
	}
}

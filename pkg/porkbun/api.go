package porkbun

import (
	"bacon/pkg/helpers"
	"encoding/json"
	"fmt"
	"io"
)

// https://porkbun.com/api/json/v3/documentation
const PING = "https://porkbun.com/api/json/v3/ping"

type baseRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type checkable interface {
	checkStatus() bool
	messageAsError() error
}

func (res *baseRes) checkStatus() bool {
	switch res.Status {
	case "SUCCESS":
		return true
	default:
		return false
	}
}

func (res *baseRes) messageAsError() error {
	return fmt.Errorf("%s", res.Message)
}

func postAndRead(domain string, body interface{}) ([]byte, error) {
	res, err := helpers.PostJson(PING, body)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}

func unmarshalAndCheckStatus(data []byte, body checkable) error {
	err := json.Unmarshal(data, &body)
	if err != nil {
		return err
	}

	if !body.checkStatus() {
		return body.messageAsError()
	}

	return nil
}

func ping(auth PorkAuth) error {
	type pingRes struct {
		baseRes
		YourIp string `json:"yourIp"`
	}

	body, err := postAndRead(PING, auth)
	if err != nil {
		return err
	}

	ping := pingRes{}
	err = unmarshalAndCheckStatus(body, &ping)
	if err != nil {
		return err
	}

	return nil
}

func deploy(auth PorkAuth, domain string, shouldCreate bool, shouldDelete bool) error {
	return fmt.Errorf("haven't implemented sync yet")
}

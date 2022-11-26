package api

import "testing"

func TestBaseResSuccess(t *testing.T) {
	res := baseRes{
		Status:  "SUCCESS",
		Message: "",
	}

	if res.checkStatus() != nil {
		t.Error("status mismatch: error when successful")
	}
}

func TestBaseResFailure(t *testing.T) {
	res := baseRes{
		Status:  "FAILURE",
		Message: "random error message",
	}

	if res.checkStatus() == nil {
		t.Error("status mismatch: no error when unsuccessful")
	}
}

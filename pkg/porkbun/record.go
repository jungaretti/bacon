package porkbun

import "bacon/pkg/client"

type porkbunRecord struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
	Notes    string `json:"notes"`
}

func ConvertToPorkbunRecord(src client.Record) (out porkbunRecord) {
	return out
}

func ConvertToClientRecord(src porkbunRecord) (out client.Record) {
	return out
}

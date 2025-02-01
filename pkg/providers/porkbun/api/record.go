package api

import "bacon/pkg/dns"

type Record struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
	Notes    string `json:"notes"`
}

func (r Record) GetName() string {
	return r.Name
}

func (r Record) GetType() string {
	return r.Type
}

func (r Record) GetTtl() string {
	return r.TTL
}

func (r Record) GetData() string {
	return r.Content
}

func (r Record) GetPriority() string {
	if r.Priority == "0" {
		return ""
	}

	return r.Priority
}

var _ dns.Record = Record{}

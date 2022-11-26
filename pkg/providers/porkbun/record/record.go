package record

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

func (r Record) GetContent() string {
	return r.Content
}

func (r Record) GetHost() string {
	return r.Name
}

func (r Record) GetTTL() string {
	return r.TTL
}

func (r Record) GetType() string {
	return r.Type
}

var _ dns.Record = Record{}

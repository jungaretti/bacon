package porkbun

import (
	"bacon/pkg/client"
	"fmt"
	"strconv"
)

type PorkbunRecord struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
	Notes    string `json:"notes"`
}

func ConvertToPorkbunRecord(src client.Record) (out PorkbunRecord) {
	out.Name = src.Host
	out.Type = src.Type
	out.Content = src.Content

	out.TTL = fmt.Sprint(src.TTL)
	out.Priority = fmt.Sprint(src.Priority)

	return out
}

func ConvertToClientRecord(src PorkbunRecord) (out client.Record, err error) {
	out.Host = src.Name
	out.Type = src.Type
	out.Content = src.Content

	if src.TTL != "" {
		ttl, err := strconv.Atoi(src.TTL)
		if err != nil {
			return out, err
		}
		out.TTL = int(ttl)
	}
	if src.Priority != "" {
		priority, err := strconv.Atoi(src.Priority)
		if err != nil {
			return out, err
		}
		out.Priority = int(priority)
	}

	return out, nil
}

func (src *PorkbunRecord) HashFuzzy() string {
	return fmt.Sprint(src.Name, src.Type, src.Content, src.TTL, src.Priority)
}

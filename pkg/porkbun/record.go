package porkbun

import (
	"bacon/pkg/client"
	"fmt"
	"strconv"
	"strings"
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

	if src.TTL > 0 {
		out.TTL = fmt.Sprint(src.TTL)
	}
	if src.Priority > 0 {
		out.Priority = fmt.Sprint(src.Priority)
	}

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

func (src PorkbunRecord) HashFuzzy() string {
	// Ignore "0"
	ttl := src.TTL
	if ttl == "0" {
		ttl = ""
	}
	priority := src.Priority
	if priority == "0" {
		priority = ""
	}

	return fmt.Sprint(src.Name, src.Type, src.Content, ttl, priority)
}

// Trims anything.domain.tld to anything
func trimDomain(name string, domain string) string {
	if name == domain {
		return ""
	}

	return strings.Replace(name, "."+domain, "", 1)
}

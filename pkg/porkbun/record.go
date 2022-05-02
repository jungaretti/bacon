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

func (left *PorkbunRecord) FuzzyCompare(right *PorkbunRecord) bool {
	// Ignore IDs and notes
	return (left.Name == right.Name &&
		left.Type == right.Type &&
		left.Content == right.Content &&
		left.TTL == right.TTL &&
		left.Priority == right.Priority)
}

func (left *PorkbunRecord) FuzzyCompareToClientRecord(right *client.Record) bool {
	// Left must be a valid int and equal to right
	if left.TTL != "" && right.TTL != 0 {
		ttl, err := strconv.Atoi(left.TTL)
		if err != nil || int(ttl) != right.TTL {
			return false
		}
	}
	if left.Priority != "" && right.Priority != 0 {
		priority, err := strconv.Atoi(left.Priority)
		if err != nil || int(priority) != right.Priority {
			return false
		}
	}

	// Ignore IDs and notes
	return (left.Name == right.Host &&
		left.Type == right.Type &&
		left.Content == right.Content)
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

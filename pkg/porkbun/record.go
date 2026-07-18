package porkbun

import "strings"

type Record struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"prio"`
	Notes    string `json:"notes"`
}

func RecordHash(r Record) string {
	return strings.Join([]string{r.Name, r.Type, r.TTL, r.Content, r.Priority, r.Notes}, "-")
}

func RecordIdentityHash(r Record) string {
	return strings.Join([]string{r.Name, r.Type}, "-")
}

func (record Record) NormalizePriority() Record {
	// Porkbun returns priority "0" for records that don't have priority
	if record.Priority == "0" {
		record.Priority = ""
	}

	return record
}

func (r Record) isIgnored() bool {
	if r.Type == "NS" {
		return true
	}
	if strings.HasPrefix(r.Name, "_acme-challenge") {
		return true
	}

	return false
}

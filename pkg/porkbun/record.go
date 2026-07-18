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
	// Treat "0" as empty for hashing purposes
	priority := r.Priority
	if priority == "0" {
		priority = ""
	}

	return strings.Join([]string{r.Name, r.Type, r.TTL, r.Content, priority, r.Notes}, "-")
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

type RecordIdentity struct {
	Name string
	Type string
}

func (r Record) Identity() RecordIdentity {
	return RecordIdentity{
		Name: r.Name,
		Type: r.Type,
	}
}

func RecordIdentityHash(identity RecordIdentity) string {
	return strings.Join([]string{identity.Name, identity.Type}, "-")
}

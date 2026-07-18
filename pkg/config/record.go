package config

import (
	"bacon/pkg/porkbun"
	"fmt"
	"strconv"
)

type Record struct {
	Name     string `yaml:"host"`
	Type     string `yaml:"type"`
	Ttl      int    `yaml:"ttl"`
	Data     string `yaml:"content"`
	Priority int    `yaml:"priority,omitempty"`
}

func (r Record) ToPorkbun() porkbun.Record {
	priority := ""
	if r.Priority != 0 {
		priority = strconv.Itoa(r.Priority)
	}

	return porkbun.Record{
		Name:     r.Name,
		Type:     r.Type,
		TTL:      strconv.Itoa(r.Ttl),
		Content:  r.Data,
		Priority: priority,
	}
}

func RecordFromPorkbun(record porkbun.Record) (Record, error) {
	ttl, err := strconv.Atoi(record.TTL)
	if err != nil {
		return Record{}, fmt.Errorf("record %v has invalid TTL: %v", record, err)
	}

	priority := 0
	if record.Priority != "" {
		priority, err = strconv.Atoi(record.Priority)
		if err != nil {
			return Record{}, fmt.Errorf("record %v has invalid priority: %v", record, err)
		}
	}

	return Record{
		Name:     record.Name,
		Type:     record.Type,
		Ttl:      ttl,
		Data:     record.Content,
		Priority: priority,
	}, nil
}

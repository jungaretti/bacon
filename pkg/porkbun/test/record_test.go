package porkbun

import (
	"bacon/pkg/client"
	"bacon/pkg/porkbun"
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConvertToPorkbun(t *testing.T) {
	record := client.Record{
		Type:    "A",
		Host:    "www.example.com",
		Content: "123.456.789.112",
		TTL:     600,
	}
	porkRecord := porkbun.ConvertToPorkbunRecord(record)

	if porkRecord.Type != record.Type {
		t.Log("expected", record.Type, "found", porkRecord.Type)
		t.Fail()
	}
	if porkRecord.Name != record.Host {
		t.Log("expected", record.Host, "found", porkRecord.Name)
		t.Fail()
	}
	if porkRecord.Content != record.Content {
		t.Log("expected", record.Content, "found", porkRecord.Content)
		t.Fail()
	}
	if porkRecord.TTL != fmt.Sprint(record.TTL) {
		t.Log("expected", fmt.Sprint(record.TTL), "found", porkRecord.TTL)
	}
	if porkRecord.Priority != fmt.Sprint(record.Priority) {
		t.Log("expected", fmt.Sprint(record.Priority), "found", porkRecord.Priority)
	}
}

func TestConvertToClientSimple(t *testing.T) {
	porkRecord := porkbun.PorkbunRecord{
		Id:      "abcxyz",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
	}
	record, err := porkbun.ConvertToClientRecord(porkRecord)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if record.Type != porkRecord.Type {
		t.Log("expected", porkRecord.Type, "found", record.Type)
		t.Fail()
	}
	if record.Host != porkRecord.Name {
		t.Log("expected", porkRecord.Name, "found", record.Host)
		t.Fail()
	}
	if record.Content != porkRecord.Content {
		t.Log("expected", porkRecord.Content, "found", record.Content)
		t.Fail()
	}
	if record.TTL != 600 {
		t.Log("expected", porkRecord.TTL, "found", 600)
		t.Fail()
	}
	if record.Priority != 0 {
		t.Log("expected", porkRecord.Priority, "found", 0)
		t.Fail()
	}
}

func TestHashFuzzy(t *testing.T) {
	one := porkbun.PorkbunRecord{
		Id:      "abc",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note2",
	}
	two := porkbun.PorkbunRecord{
		Id:      "xyz",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note2",
	}
	three := porkbun.PorkbunRecord{
		Id:      "xyz",
		Name:    "www.example.org",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note2",
	}

	if one.HashFuzzy() != two.HashFuzzy() {
		t.Log("left", one)
		t.Log("right", two)
		t.Fail()
	}
	if two.HashFuzzy() == three.HashFuzzy() {
		t.Log("left", two)
		t.Log("right", three)
		t.Fail()
	}
}

func TestComparePorkbunResponse(t *testing.T) {
	const config = `type: CNAME
host: www.jungaretti.com
content: green-forest-08be42d1e.azurestaticapps.net
ttl: 600`
	const response = `{"id": "188254337","name": "www.jungaretti.com","type": "CNAME","content": "green-forest-08be42d1e.azurestaticapps.net","ttl": "600","prio": "0"}`

	record := client.Record{}
	err := yaml.Unmarshal([]byte(config), &record)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	porkRecord := porkbun.PorkbunRecord{}
	err = json.Unmarshal([]byte(response), &porkRecord)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	convertedRecord := porkbun.ConvertToPorkbunRecord(record)
	if porkRecord.HashFuzzy() != convertedRecord.HashFuzzy() {
		t.Log(porkRecord.HashFuzzy())
		t.Log(convertedRecord.HashFuzzy())
		t.Fail()
	}
}

func TestComparePorkbunResponseNullPriority(t *testing.T) {
	const config = `type: CNAME
host: www.jungaretti.com
content: green-forest-08be42d1e.azurestaticapps.net
ttl: 600`
	const response = `{"id": "188254337","name": "www.jungaretti.com","type": "CNAME","content": "green-forest-08be42d1e.azurestaticapps.net","ttl": "600","prio": null}`

	record := client.Record{}
	err := yaml.Unmarshal([]byte(config), &record)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	porkRecord := porkbun.PorkbunRecord{}
	err = json.Unmarshal([]byte(response), &porkRecord)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	convertedRecord := porkbun.ConvertToPorkbunRecord(record)
	if porkRecord.HashFuzzy() != convertedRecord.HashFuzzy() {
		t.Log("one", porkRecord.HashFuzzy())
		t.Log("two", convertedRecord.HashFuzzy())
		t.Error()
	}
}

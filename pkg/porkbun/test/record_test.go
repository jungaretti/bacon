package porkbun

import (
	"bacon/pkg/client"
	"bacon/pkg/porkbun"
	"fmt"
	"testing"
)

func TestConvertToPorkbunRecord(t *testing.T) {
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

func TestConvertToClientRecord(t *testing.T) {
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

	if record.Host != porkRecord.Name {
		t.Log("expected", porkRecord.Name, "found", record.Host)
		t.Fail()
	}
	if record.Type != porkRecord.Type {
		t.Log("expected", porkRecord.Type, "found", record.Type)
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

func TestFuzzyCompare(t *testing.T) {
	left := porkbun.PorkbunRecord{
		Id:      "abc",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note1",
	}
	right := porkbun.PorkbunRecord{
		Id:      "xyz",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note2",
	}

	if !left.FuzzyCompare(&right) {
		t.Log("left", left)
		t.Log("right", right)
		t.Fail()
	}
}

func TestFuzzyCompareDistinct(t *testing.T) {
	left := porkbun.PorkbunRecord{
		Id:      "abc",
		Name:    "mail.example.com",
		Type:    "A",
		Content: "112.123.789.456",
		TTL:     "600",
		Notes:   "note3",
	}
	right := porkbun.PorkbunRecord{
		Id:      "xyz",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note2",
	}

	if left.FuzzyCompare(&right) {
		t.Log("left", left)
		t.Log("right", right)
		t.Fail()
	}
}

func TestFuzzyCompareToClientRecord(t *testing.T) {
	left := porkbun.PorkbunRecord{
		Id:      "abc",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note3",
	}
	right := client.Record{
		Host:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     600,
	}

	if !left.FuzzyCompareToClientRecord(&right) {
		t.Log("left", left)
		t.Log("right", right)
		t.Fail()
	}
}

func TestFuzzyCompareToClientRecordDistinct(t *testing.T) {
	left := porkbun.PorkbunRecord{
		Id:      "abc",
		Name:    "mail.example.com",
		Type:    "A",
		Content: "112.123.789.456",
		TTL:     "600",
		Notes:   "note3",
	}
	right := client.Record{
		Host:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     600,
	}

	if left.FuzzyCompareToClientRecord(&right) {
		t.Log("left", left)
		t.Log("right", right)
		t.Fail()
	}
}

func TestFuzzyHash(t *testing.T) {
	left := porkbun.PorkbunRecord{
		Id:      "abc",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note1",
	}
	right := porkbun.PorkbunRecord{
		Id:      "xyz",
		Name:    "www.example.com",
		Type:    "A",
		Content: "123.456.789.112",
		TTL:     "600",
		Notes:   "note2",
	}

	if left.FuzzyHash() != right.FuzzyHash() {
		t.Log("left", left)
		t.Log("right", right)
		t.Fail()
	}
}

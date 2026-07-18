package porkbun

import "testing"

func TestDiffRecordsChangedName(t *testing.T) {
	from := []Record{
		{
			Id:      "123",
			Name:    "www.bacondemo.com",
			Type:    "A",
			TTL:     "600",
			Content: "1.2.3.4",
		},
	}

	to := []Record{
		{
			Name:    "api.bacondemo.com",
			Type:    "A",
			TTL:     "600",
			Content: "1.2.3.4",
		},
	}

	added, removed, edited, unchanged := DiffRecords(from, to)

	if len(added) != 1 {
		t.Error("expected 1 added, got", len(added))
	}

	if len(removed) != 1 {
		t.Error("expected 1 removed, got", len(removed))
	}

	if len(edited) != 0 {
		t.Error("expected no edited, got", len(edited))
	}

	if len(unchanged) != 0 {
		t.Error("expected no unchanged, got", len(unchanged))
	}
}

func TestDiffRecordsChangedType(t *testing.T) {
	from := []Record{
		{
			Id:      "123",
			Name:    "www.bacondemo.com",
			Type:    "A",
			TTL:     "600",
			Content: "1.2.3.4",
		},
	}

	to := []Record{
		{
			Name:    "www.bacondemo.com",
			Type:    "CNAME",
			TTL:     "600",
			Content: "somewhere.else.com",
		},
	}

	added, removed, edited, unchanged := DiffRecords(from, to)

	if len(added) != 1 {
		t.Error("expected 1 added, got", len(added))
	}

	if len(removed) != 1 {
		t.Error("expected 1 removed, got", len(removed))
	}

	if len(edited) != 0 {
		t.Error("expected no edited, got", len(edited))
	}

	if len(unchanged) != 0 {
		t.Error("expected no unchanged, got", len(unchanged))
	}
}

func TestDiffRecordsChangedContent(t *testing.T) {
	from := []Record{
		{
			Id:      "123",
			Name:    "www.bacondemo.com",
			Type:    "A",
			TTL:     "600",
			Content: "1.2.3.4",
		},
	}

	to := []Record{
		{
			Name:    "www.bacondemo.com",
			Type:    "A",
			TTL:     "600",
			Content: "5.6.7.8",
		},
	}

	added, removed, edited, unchanged := DiffRecords(from, to)

	if len(added) != 0 {
		t.Error("expected no added, got", len(added))
	}

	if len(removed) != 0 {
		t.Error("expected no removed, got", len(removed))
	}

	if len(edited) != 1 {
		t.Error("expected 1 edited, got", len(edited))
	}
	if edited[0].Id != "123" {
		t.Error("expected edit to keep ID 123, got", edited[0].Id)
	}
	if edited[0].Content != "5.6.7.8" {
		t.Error("expected edit to carry new content, got", edited[0].Content)
	}

	if len(unchanged) != 0 {
		t.Error("expected no unchanged, got", len(unchanged))
	}
}

func TestDiffRecordsChangedTTL(t *testing.T) {
	from := []Record{
		{
			Id:      "123",
			Name:    "www.bacondemo.com",
			Type:    "A",
			TTL:     "600",
			Content: "1.2.3.4",
		},
	}

	to := []Record{
		{
			Name:    "www.bacondemo.com",
			Type:    "A",
			TTL:     "3600",
			Content: "1.2.3.4",
		},
	}

	added, removed, edited, unchanged := DiffRecords(from, to)

	if len(added) != 0 {
		t.Error("expected no added, got", len(added))
	}

	if len(removed) != 0 {
		t.Error("expected no removed, got", len(removed))
	}

	if len(edited) != 1 {
		t.Error("expected 1 edited, got", len(edited))
	}
	if edited[0].Id != "123" {
		t.Error("expected edit to keep ID 123, got", edited[0].Id)
	}
	if edited[0].TTL != "3600" {
		t.Error("expected edit to carry new TTL, got", edited[0].TTL)
	}

	if len(unchanged) != 0 {
		t.Error("expected no unchanged, got", len(unchanged))
	}
}

func TestDiffRecordsChangedPriority(t *testing.T) {
	from := []Record{
		{
			Id:       "123",
			Name:     "bacondemo.com",
			Type:     "MX",
			TTL:      "600",
			Priority: "10",
			Content:  "mx.bacondemo.com",
		},
		{
			Id:       "456",
			Name:     "bacondemo.com",
			Type:     "MX",
			TTL:      "600",
			Priority: "20",
			Content:  "mx2.bacondemo.com",
		},
	}

	to := []Record{
		{
			Name:     "bacondemo.com",
			Type:     "MX",
			TTL:      "600",
			Priority: "10",
			Content:  "mx.bacondemo.com",
		},
		{
			Name:     "bacondemo.com",
			Type:     "MX",
			TTL:      "600",
			Priority: "30",
			Content:  "mx2.bacondemo.com",
		},
	}

	added, removed, edited, unchanged := DiffRecords(from, to)

	if len(added) != 0 {
		t.Error("expected no added, got", len(added))
	}

	if len(removed) != 0 {
		t.Error("expected no removed, got", len(removed))
	}

	if len(edited) != 1 {
		t.Error("expected 1 edited, got", len(edited))
	}
	if edited[0].Id != "456" {
		t.Error("expected edit to keep ID 456, got", edited[0].Id)
	}
	if edited[0].Priority != "30" {
		t.Error("expected edit to carry new priority, got", edited[0].Priority)
	}

	if len(unchanged) != 1 {
		t.Error("expected 1 unchanged, got", len(unchanged))
	}
}

func TestDiffRecordsChangedNotes(t *testing.T) {
	from := []Record{
		{
			Id:       "123",
			Name:     "bacondemo.com",
			Type:     "MX",
			TTL:      "600",
			Priority: "10",
			Content:  "mx.bacondemo.com",
			Notes:    "Fastmail",
		},
	}

	to := []Record{
		{
			Name:     "bacondemo.com",
			Type:     "MX",
			TTL:      "600",
			Priority: "10",
			Content:  "mx.bacondemo.com",
			Notes:    "Exchange",
		},
	}

	added, removed, edited, unchanged := DiffRecords(from, to)

	if len(added) != 0 {
		t.Error("expected no added, got", len(added))
	}

	if len(removed) != 0 {
		t.Error("expected no removed, got", len(removed))
	}

	if len(edited) != 1 {
		t.Error("expected 1 edited, got", len(edited))
	}

	if edited[0].Id != "123" {
		t.Error("expected edit to keep ID 123, got", edited[0].Id)
	}
	if edited[0].Notes != "Exchange" {
		t.Error("expected edit to carry new notes, got", edited[0].Notes)
	}

	if len(unchanged) != 0 {
		t.Error("expected no unchanged, got", len(unchanged))
	}
}

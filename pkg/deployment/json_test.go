package deployment

import (
	"bacon/pkg/porkbun"
	"bytes"
	"encoding/json"
	"testing"
)

func TestJsonReporterFinish(t *testing.T) {
	results := []Operation{
		{Action: Delete, DryRun: false, Record: porkbun.Record{Name: "old.example.com", Type: "A", Content: "1.2.3.4", TTL: "600"}},
		{Action: Update, DryRun: true, Record: porkbun.Record{Name: "api.example.com", Type: "A", Content: "42.42.42.42", TTL: "600", Notes: "API endpoint"}},
		{Action: Create, DryRun: true, Record: porkbun.Record{Name: "example.com", Type: "MX", Content: "smtp.example.com", TTL: "600", Priority: "10"}},
		{Action: Keep, Record: porkbun.Record{Name: "example.com", Type: "ALIAS", Content: "host.example.com", TTL: "600"}},
	}

	buffer := bytes.Buffer{}
	reporter := jsonReporter{writer: &buffer}
	err := reporter.Report(results)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsed := struct {
		Summary map[string]int   `json:"summary"`
		Records []map[string]any `json:"records"`
	}{}
	err = json.Unmarshal(buffer.Bytes(), &parsed)
	if err != nil {
		t.Fatalf("expected valid JSON, got %v", err)
	}

	expectedSummary := map[string]int{
		"deleted":     1,
		"wouldDelete": 0,
		"updated":     0,
		"wouldUpdate": 1,
		"created":     0,
		"wouldCreate": 1,
		"kept":        1,
	}
	for key, count := range expectedSummary {
		if parsed.Summary[key] != count {
			t.Errorf("expected summary %v to be %v, got %v", key, count, parsed.Summary[key])
		}
	}

	if len(parsed.Records) != 4 {
		t.Fatalf("expected 4 records, got %v", len(parsed.Records))
	}
	if parsed.Records[0]["action"] != "delete" {
		t.Errorf("expected first action to be delete, got %v", parsed.Records[0]["action"])
	}
	if parsed.Records[0]["dryRun"] != false {
		t.Errorf("expected first record to not be a dry run, got %v", parsed.Records[0]["dryRun"])
	}
	if parsed.Records[1]["dryRun"] != true {
		t.Errorf("expected second record to be a dry run, got %v", parsed.Records[1]["dryRun"])
	}

	record := parsed.Records[2]["record"].(map[string]any)
	if record["priority"] != "10" {
		t.Errorf("expected priority key with value 10, got %v", record["priority"])
	}
	if record["name"] != "example.com" {
		t.Errorf("expected name example.com, got %v", record["name"])
	}
}

package operation

import (
	"bacon/pkg/porkbun"
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestTableReporter(t *testing.T) {
	buffer := &bytes.Buffer{}
	reporter := NewTableReporter(buffer)

	reporter.Report(RecordOperationResult{
		Operation: RecordOperation{
			Action: Create,
			Record: porkbun.Record{
				Name:    "bacontest42.com",
				Type:    "ALIAS",
				Content: "pixie.porkbun.com",
				TTL:     "600",
			},
		},
	})
	reporter.Report(RecordOperationResult{
		Operation: RecordOperation{
			Action: Delete,
			Record: porkbun.Record{
				Name:    "old.bacontest42.com",
				Type:    "A",
				Content: "1.2.3.4",
				TTL:     "600",
			},
		},
		Err: errors.New("api error"),
	})
	reporter.Report(RecordOperationResult{
		Operation: RecordOperation{
			Action: Skip,
			Record: porkbun.Record{
				Name:    "www.bacontest42.com",
				Type:    "CNAME",
				Content: "pixie.porkbun.com",
				TTL:     "600",
			},
		},
	})
	reporter.Finish(Summary{Created: 1, Unchanged: 1, Failed: 1})

	output := buffer.String()
	findRow := func(substring string) string {
		for _, line := range strings.Split(output, "\n") {
			if strings.Contains(line, substring) {
				return line
			}
		}
		t.Fatal("no row contains", substring)
		return ""
	}

	rows := map[string][]string{
		"ALIAS":               {"create", "succeeded"},
		"old.bacontest42.com": {"delete", "failed"},
		"www.bacontest42.com": {"unchanged"},
	}
	for row, marks := range rows {
		line := findRow(row)
		for _, mark := range marks {
			if !strings.Contains(line, mark) {
				t.Error("expected row", line, "to contain", mark)
			}
		}
	}

	for _, expected := range []string{
		"ACTION",
		"error: api error",
		"1 created, 1 unchanged, 1 failed",
	} {
		if !strings.Contains(output, expected) {
			t.Error("expected output to contain", expected, "actual", output)
		}
	}
}

func TestTableReporterTruncatesLongValues(t *testing.T) {
	buffer := &bytes.Buffer{}
	reporter := NewTableReporter(buffer)

	reporter.Report(RecordOperationResult{
		Operation: RecordOperation{
			Action: Create,
			Record: porkbun.Record{
				Name:    "a-very-long-subdomain-name.bacontest42.com",
				Type:    "TXT",
				Content: "in1-smtp.messagingengine.com.in1-smtp.messagingengine.com",
				TTL:     "600",
			},
		},
	})

	output := buffer.String()
	if !strings.Contains(output, "a-very-long-subdomain...") {
		t.Error("expected output to contain truncated host, actual", output)
	}
	if !strings.Contains(output, "in1-smtp.messagingengine.co...") {
		t.Error("expected output to contain truncated content, actual", output)
	}
	if strings.Contains(output, "a-very-long-subdomain-name.bacontest42.com") {
		t.Error("expected output to not contain full host, actual", output)
	}
}

func TestJSONReporter(t *testing.T) {
	buffer := &bytes.Buffer{}
	reporter := NewJSONReporter(buffer)

	reporter.Report(RecordOperationResult{
		Operation: RecordOperation{
			Action: Delete,
			Record: porkbun.Record{
				Name:    "old.bacontest42.com",
				Type:    "A",
				Content: "1.2.3.4",
				TTL:     "600",
			},
		},
		Err: errors.New("api error"),
	})
	reporter.Finish(Summary{Failed: 1})

	document := struct {
		Results []struct {
			Action string         `json:"action"`
			Record porkbun.Record `json:"record"`
			Status string         `json:"status"`
			Error  string         `json:"error"`
		} `json:"results"`
		Summary Summary `json:"summary"`
	}{}

	err := json.Unmarshal(buffer.Bytes(), &document)
	if err != nil {
		t.Fatal("couldn't unmarshal output:", err)
	}

	if len(document.Results) != 1 {
		t.Fatal("expected", 1, "actual", len(document.Results))
	}
	if document.Results[0].Action != string(Delete) {
		t.Error("expected", Delete, "actual", document.Results[0].Action)
	}
	if document.Results[0].Status != string(Failed) {
		t.Error("expected", Failed, "actual", document.Results[0].Status)
	}
	if document.Results[0].Error != "api error" {
		t.Error("expected", "api error", "actual", document.Results[0].Error)
	}
	if document.Results[0].Record.Name != "old.bacontest42.com" {
		t.Error("expected", "old.bacontest42.com", "actual", document.Results[0].Record.Name)
	}
	if document.Summary.Failed != 1 {
		t.Error("expected", 1, "actual", document.Summary.Failed)
	}
}

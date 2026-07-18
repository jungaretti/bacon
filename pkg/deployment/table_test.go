package deployment

import (
	"bytes"
	"testing"
)

func TestTableReporterSummary(t *testing.T) {
	results := []Operation{
		{Action: Delete, DryRun: false},
		{Action: Delete, DryRun: false},
		{Action: Update, DryRun: true},
		{Action: Keep},
	}

	buffer := bytes.Buffer{}
	reporter := tableReporter{writer: &buffer}
	err := reporter.Report(results)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Summary: Deleted 2 records, would update 1 record, kept 1 record\n"
	if buffer.String() != expected {
		t.Errorf("expected %q, got %q", expected, buffer.String())
	}
}

func TestTableReporterSummaryEmpty(t *testing.T) {
	buffer := bytes.Buffer{}
	reporter := tableReporter{writer: &buffer}
	err := reporter.Report([]Operation{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Summary: no records\n"
	if buffer.String() != expected {
		t.Errorf("expected %q, got %q", expected, buffer.String())
	}
}

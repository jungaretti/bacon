package operation

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Summary struct {
	Created        int `json:"created"`
	Deleted        int `json:"deleted"`
	PlannedCreates int `json:"plannedCreates"`
	PlannedDeletes int `json:"plannedDeletes"`
	Unchanged      int `json:"unchanged"`
	Failed         int `json:"failed"`
}

func (s Summary) String() string {
	var parts []string
	appendCount := func(count int, singular, plural string) {
		if count == 1 {
			parts = append(parts, fmt.Sprintf("%d %s", count, singular))
		} else if count > 1 {
			parts = append(parts, fmt.Sprintf("%d %s", count, plural))
		}
	}

	appendCount(s.Created, "created", "created")
	appendCount(s.Deleted, "deleted", "deleted")
	appendCount(s.PlannedCreates, "create planned", "creates planned")
	appendCount(s.PlannedDeletes, "delete planned", "deletes planned")
	appendCount(s.Unchanged, "unchanged", "unchanged")
	appendCount(s.Failed, "failed", "failed")

	if len(parts) == 0 {
		return "no records"
	}
	return strings.Join(parts, ", ")
}

type Reporter interface {
	Report(result RecordOperationResult)
	Finish(s Summary) error
}

// Collects results and serializes them to a single JSON object after all operations have completed.
type jsonReporter struct {
	w       io.Writer
	results []RecordOperationResult
}

// NewJSONReporter returns a Reporter that collects results and serializes them to a single JSON object after all operations have completed.
func NewJSONReporter(w io.Writer) Reporter {
	return &jsonReporter{w: w, results: []RecordOperationResult{}}
}

func (j *jsonReporter) Report(result RecordOperationResult) {
	j.results = append(j.results, result)
}

func (j *jsonReporter) Finish(summary Summary) error {
	document := struct {
		Summary Summary                 `json:"summary"`
		Results []RecordOperationResult `json:"results"`
	}{
		Results: j.results,
		Summary: summary,
	}

	encoder := json.NewEncoder(j.w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(document)
}

// Streams results in a tabular format as they complete.
type tableReporter struct {
	w           io.Writer
	wroteHeader bool
}

// NewTableReporter returns a Reporter that streams results in a tabular format as they complete.
func NewTableReporter(w io.Writer) Reporter {
	return &tableReporter{w: w}
}

func (t *tableReporter) Report(result RecordOperationResult) {
	if !t.wroteHeader {
		t.printRow("ACTION", "STATUS", "HOST", "TYPE", "CONTENT", "TTL", "PRIO")
		t.printSeparator()
		t.wroteHeader = true
	}

	record := result.Operation.Record
	priority := record.Priority
	if priority == "0" {
		priority = ""
	}

	t.printRow(string(result.Operation.Action), string(result.Status()), record.Name, record.Type, record.Content, record.TTL, priority)
	if result.Err != nil {
		fmt.Fprintln(t.w, "  error:", result.Err)
	}
}

func (t *tableReporter) Finish(summary Summary) error {
	_, err := fmt.Fprintln(t.w, "\n"+summary.String())
	return err
}

const (
	actionWidth  = 6
	statusWidth  = 9
	hostWidth    = 24
	typeWidth    = 6
	contentWidth = 30
	ttlWidth     = 5
)

func (t *tableReporter) printRow(action, status, host, recordType, content, ttl, priority string) {
	row := fmt.Sprintf("%-*s  %-*s  %-*s  %-*s  %-*s  %-*s  %s",
		actionWidth, truncate(action, actionWidth),
		statusWidth, truncate(status, statusWidth),
		hostWidth, truncate(host, hostWidth),
		typeWidth, truncate(recordType, typeWidth),
		contentWidth, truncate(content, contentWidth),
		ttlWidth, truncate(ttl, ttlWidth),
		priority)
	fmt.Fprintln(t.w, strings.TrimRight(row, " "))
}

func (t *tableReporter) printSeparator() {
	dashes := func(width int) string { return strings.Repeat("-", width) }
	t.printRow(dashes(actionWidth), dashes(statusWidth), dashes(hostWidth), dashes(typeWidth), dashes(contentWidth), dashes(ttlWidth), dashes(len("PRIO")))
}

// Shortens a value to fit within width, marking cut values with "..."
func truncate(value string, width int) string {
	if len(value) <= width {
		return value
	}
	return value[:width-3] + "..."
}

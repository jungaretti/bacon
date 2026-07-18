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

func (summary Summary) String() string {
	var parts []string
	appendCount := func(count int, singular, plural string) {
		if count == 1 {
			parts = append(parts, fmt.Sprintf("%d %s", count, singular))
		} else if count > 1 {
			parts = append(parts, fmt.Sprintf("%d %s", count, plural))
		}
	}

	appendCount(summary.Created, "created", "created")
	appendCount(summary.Deleted, "deleted", "deleted")
	appendCount(summary.PlannedCreates, "create planned", "creates planned")
	appendCount(summary.PlannedDeletes, "delete planned", "deletes planned")
	appendCount(summary.Unchanged, "unchanged", "unchanged")
	appendCount(summary.Failed, "failed", "failed")

	if len(parts) == 0 {
		return "no records"
	}
	return strings.Join(parts, ", ")
}

type Reporter interface {
	Report(result RecordOperationResult)
	Finish(summary Summary) error
}

type jsonReporter struct {
	writer  io.Writer
	results []RecordOperationResult
}

// NewJSONReporter returns a Reporter that collects results and serializes
// them to a single JSON object after all operations have completed.
func NewJSONReporter(writer io.Writer) Reporter {
	return &jsonReporter{writer: writer, results: []RecordOperationResult{}}
}

func (reporter *jsonReporter) Report(result RecordOperationResult) {
	reporter.results = append(reporter.results, result)
}

func (reporter *jsonReporter) Finish(summary Summary) error {
	document := struct {
		Summary Summary                 `json:"summary"`
		Results []RecordOperationResult `json:"results"`
	}{
		Summary: summary,
		Results: reporter.results,
	}

	encoder := json.NewEncoder(reporter.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(document)
}

var tableColumns = []struct {
	header string
	width  int
}{
	{"ACTION", 6},
	{"STATUS", 9},
	{"HOST", 24},
	{"TYPE", 6},
	{"CONTENT", 30},
	{"TTL", 5},
	{"PRIO", 4},
}

type tableReporter struct {
	writer      io.Writer
	wroteHeader bool
}

// NewTableReporter returns a Reporter that streams results in a tabular
// format as they complete.
func NewTableReporter(writer io.Writer) Reporter {
	return &tableReporter{writer: writer}
}

func (reporter *tableReporter) Report(result RecordOperationResult) {
	if !reporter.wroteHeader {
		reporter.printHeader()
		reporter.wroteHeader = true
	}

	record := result.Operation.Record
	priority := record.Priority
	if priority == "0" {
		priority = ""
	}

	reporter.printRow(actionLabel(result.Operation.Action), string(result.Status()), record.Name, record.Type, record.Content, record.TTL, priority)
	if result.Err != nil {
		fmt.Fprintln(reporter.writer, "  error:", result.Err)
	}
}

func (reporter *tableReporter) Finish(summary Summary) error {
	_, err := fmt.Fprintln(reporter.writer, "\n"+summary.String())
	return err
}

func (reporter *tableReporter) printHeader() {
	headers := make([]string, len(tableColumns))
	separators := make([]string, len(tableColumns))
	for i, column := range tableColumns {
		headers[i] = column.header
		separators[i] = strings.Repeat("-", column.width)
	}

	reporter.printRow(headers...)
	reporter.printRow(separators...)
}

func (reporter *tableReporter) printRow(values ...string) {
	cells := make([]string, len(values))
	for i, value := range values {
		width := tableColumns[i].width
		cells[i] = fmt.Sprintf("%-*s", width, truncate(value, width))
	}

	row := strings.Join(cells, "  ")
	fmt.Fprintln(reporter.writer, strings.TrimRight(row, " "))
}

// Label for the table's ACTION column. Skip becomes blank because it doesn't do anything.
func actionLabel(action Action) string {
	if action == Skip {
		return ""
	}
	return string(action)
}

// Shortens a value to fit within width, marking cut values with "..."
func truncate(value string, width int) string {
	if len(value) <= width {
		return value
	}
	return value[:width-3] + "..."
}

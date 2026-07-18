package deployment

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const (
	createSymbol = "+"
	updateSymbol = "~"
	deleteSymbol = "-"
	keepSymbol   = "="
)

type tableReporter struct {
	writer io.Writer
}

func (reporter *tableReporter) Preview(results []Operation) error {
	output := table.New("", "NAME", "TYPE", "CONTENT", "TTL", "PRI", "NOTES").
		WithWriter(reporter.writer).
		WithHeaderFormatter(color.New(color.Underline).SprintfFunc()).
		WithWidthFunc(visibleWidth)

	for _, result := range results {
		record := result.Record
		output.AddRow(actionSymbol(result.Action), record.Name, record.Type, record.Content, record.TTL, record.Priority, record.Notes)
	}

	output.Print()
	_, err := fmt.Fprintln(reporter.writer)
	return err
}

func (reporter *tableReporter) Report(results []Operation) error {
	counts := map[Action]int{}
	dryRun := map[Action]bool{}
	for _, result := range results {
		counts[result.Action]++
		dryRun[result.Action] = result.DryRun
	}

	clauses := []string{}
	if counts[Delete] > 0 {
		clauses = append(clauses, summarizeAction(dryRun[Delete], "would delete", "deleted", counts[Delete]))
	}
	if counts[Update] > 0 {
		clauses = append(clauses, summarizeAction(dryRun[Update], "would update", "updated", counts[Update]))
	}
	if counts[Create] > 0 {
		clauses = append(clauses, summarizeAction(dryRun[Create], "would create", "created", counts[Create]))
	}
	if counts[Keep] > 0 {
		clauses = append(clauses, "kept "+countRecords(counts[Keep]))
	}

	if len(clauses) == 0 {
		_, err := fmt.Fprintln(reporter.writer, "Summary: no records")
		return err
	}

	summary := strings.Join(clauses, ", ")
	_, err := fmt.Fprintln(reporter.writer, "Summary:", strings.ToUpper(summary[:1])+summary[1:])
	return err
}

func actionSymbol(action Action) string {
	switch action {
	case Delete:
		return color.RedString(deleteSymbol)
	case Update:
		return color.YellowString(updateSymbol)
	case Create:
		return color.GreenString(createSymbol)
	}

	return keepSymbol
}

var ansiEscapes = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Colored cells must not count their ANSI escapes toward column width.
func visibleWidth(text string) int {
	return utf8.RuneCountInString(ansiEscapes.ReplaceAllString(text, ""))
}

func summarizeAction(dryRun bool, plannedVerb string, appliedVerb string, count int) string {
	verb := appliedVerb
	if dryRun {
		verb = plannedVerb
	}

	return verb + " " + countRecords(count)
}

func countRecords(count int) string {
	if count == 1 {
		return "1 record"
	}

	return fmt.Sprintf("%d records", count)
}

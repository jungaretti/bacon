package deployment

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

type TableFormatter struct{}

var _ Formatter = TableFormatter{}

func (formatter TableFormatter) FormatStart(domain string, dryRun bool) string {
	if dryRun {
		return fmt.Sprintf("Dry-run deployment of %v:\n", domain)
	}
	return fmt.Sprintf("Deploying %v...\n", domain)
}

func (formatter TableFormatter) FormatResult(deploymentResult DeploymentResult) string {
	results := deploymentResult.Results

	hasError := false
	for _, result := range results {
		if result.Error != "" {
			hasError = true
			break
		}
	}

	var builder strings.Builder
	if len(results) > 0 {
		writer := tabwriter.NewWriter(&builder, 0, 4, 1, ' ', 0)
		if hasError {
			fmt.Fprintln(writer, "OP\tNAME\tTYPE\tCONTENT\tTTL\tPRI\tNOTES\tERROR")
		} else {
			fmt.Fprintln(writer, "OP\tNAME\tTYPE\tCONTENT\tTTL\tPRI\tNOTES")
		}
		for _, result := range results {
			record := result.Record
			priority := ""
			if record.Priority != 0 {
				priority = fmt.Sprint(record.Priority)
			}
			if hasError {
				fmt.Fprintf(writer, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", operationSymbol(result.Type), record.Name, record.Type, record.Data, record.Ttl, priority, record.Notes, result.Error)
			} else {
				fmt.Fprintf(writer, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", operationSymbol(result.Type), record.Name, record.Type, record.Data, record.Ttl, priority, record.Notes)
			}
		}
		writer.Flush()
		fmt.Fprintln(&builder)
	}

	fmt.Fprint(&builder, summaryLines(deploymentResult.Summary))
	return builder.String()
}

func operationSymbol(operationType OperationType) string {
	switch operationType {
	case Create:
		return "+"
	case Update:
		return "~"
	case Delete:
		return "-"
	case Keep:
		return "="
	}
	return "?"
}

func summaryLines(summary DeploymentSummary) string {
	type verbGroup struct {
		operationType OperationType
		plannedVerb   string
		executedVerb  string
	}
	groups := []verbGroup{
		{Delete, "delete", "deleted"},
		{Update, "update", "updated"},
		{Create, "create", "created"},
		{Keep, "keep", "kept"},
	}

	var builder strings.Builder
	for _, group := range groups {
		if count := summary.OperationCounts[group.operationType]; count > 0 {
			phrase := group.executedVerb + " " + pluralizeRecords(count)
			if summary.DryRun && group.operationType != Keep {
				phrase = "would " + group.plannedVerb + " " + pluralizeRecords(count)
			}
			fmt.Fprintf(&builder, "%v %v\n", operationSymbol(group.operationType), phrase)
		}
		if count := summary.FailureCounts[group.operationType]; count > 0 {
			fmt.Fprintf(&builder, "%v failed to %v %v\n", operationSymbol(group.operationType), group.plannedVerb, pluralizeRecords(count))
		}
	}

	if builder.Len() == 0 {
		return "No records\n"
	}
	return builder.String()
}

func pluralizeRecords(count int) string {
	if count == 1 {
		return "1 record"
	}
	return fmt.Sprintf("%d records", count)
}

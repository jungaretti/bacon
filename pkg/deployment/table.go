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
			if hasError {
				fmt.Fprintf(writer, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", operationSymbol(result.Type), record.Name, record.Type, record.Content, record.TTL, record.Priority, record.Notes, result.Error)
			} else {
				fmt.Fprintf(writer, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", operationSymbol(result.Type), record.Name, record.Type, record.Content, record.TTL, record.Priority, record.Notes)
			}
		}
		writer.Flush()
		fmt.Fprintln(&builder)
	}

	fmt.Fprintln(&builder, summarySentence(deploymentResult.Summary))
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

func summarySentence(summary DeploymentSummary) string {
	changeGroups := []struct {
		operationType OperationType
		plannedVerb   string
		executedVerb  string
	}{
		{Delete, "delete", "deleted"},
		{Update, "update", "updated"},
		{Create, "create", "created"},
	}

	var changePhrases []string
	for _, group := range changeGroups {
		count := summary.OperationCounts[group.operationType]
		if count == 0 {
			continue
		}
		verb := group.executedVerb
		if summary.DryRun {
			verb = group.plannedVerb
		}
		changePhrases = append(changePhrases, verb+" "+pluralizeRecords(count))
	}

	keptPhrase := ""
	if keptCount := summary.OperationCounts[Keep]; keptCount > 0 {
		keptPhrase = "kept " + pluralizeRecords(keptCount)
	}

	if len(changePhrases) == 0 && keptPhrase == "" {
		return "No records"
	}
	if len(changePhrases) == 0 {
		return capitalize(keptPhrase)
	}

	if summary.DryRun {
		sentence := "Would " + strings.Join(changePhrases, ", ")
		if keptPhrase != "" {
			sentence += "; " + keptPhrase
		}
		return sentence
	}

	phrases := changePhrases
	if keptPhrase != "" {
		phrases = append(phrases, keptPhrase)
	}
	return capitalize(strings.Join(phrases, ", "))
}

func pluralizeRecords(count int) string {
	if count == 1 {
		return "1 record"
	}
	return fmt.Sprintf("%d records", count)
}

func capitalize(sentence string) string {
	if sentence == "" {
		return sentence
	}
	return strings.ToUpper(sentence[:1]) + sentence[1:]
}

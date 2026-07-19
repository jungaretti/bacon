package deployment

import (
	"fmt"
	"strings"
)

type TextFormatter struct{}

var _ Formatter = TextFormatter{}

const (
	createSymbol = "+"
	updateSymbol = "~"
	deleteSymbol = "-"
	keepSymbol   = "="
)

func (formatter TextFormatter) Format(deploymentResult DeploymentResult) string {
	results := deploymentResult.Results

	executed := false
	for _, result := range results {
		if result.Status != Planned {
			executed = true
			break
		}
	}

	var builder strings.Builder
	formatSection(&builder, results, Delete, deleteSymbol, "Deleting", "Would delete", executed)
	formatSection(&builder, results, Update, updateSymbol, "Updating", "Would update", executed)
	formatSection(&builder, results, Create, createSymbol, "Creating", "Would create", executed)

	fmt.Fprintln(&builder, "Keeping", countByType(results, Keep), "records:")
	for _, result := range results {
		if result.Type == Keep {
			fmt.Fprintln(&builder, keepSymbol, result.Record)
		}
	}

	if executed {
		fmt.Fprintln(&builder, "Deployment complete!")
	} else {
		fmt.Fprintln(&builder, "Mock deployment complete")
	}
	return builder.String()
}

func formatSection(builder *strings.Builder, results []OperationResult, operationType OperationType, symbol, executedVerb, plannedVerb string, executed bool) {
	count := countByType(results, operationType)
	if executed {
		fmt.Fprintln(builder, executedVerb, count, "records...")
	} else {
		fmt.Fprintln(builder, plannedVerb, count, "records:")
	}

	for _, result := range results {
		if result.Type != operationType {
			continue
		}
		if result.Status == Failure {
			fmt.Fprintf(builder, "%v %v (failed: %v)\n", symbol, result.Record, result.Error)
		} else {
			fmt.Fprintln(builder, symbol, result.Record)
		}
	}
}

func countByType(results []OperationResult, operationType OperationType) int {
	count := 0
	for _, result := range results {
		if result.Type == operationType {
			count++
		}
	}
	return count
}

package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/deployment"
	"bacon/pkg/porkbun"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	createSymbol = "+"
	updateSymbol = "~"
	deleteSymbol = "-"
	keepSymbol   = "="
)

func newDeployCmd(client *porkbun.Client) *cobra.Command {
	var shouldCreate bool
	var shouldDelete bool
	var shouldUpdate bool

	deploy := &cobra.Command{
		Use:   "deploy <config-file>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by updating records in place when
possible, then deleting old records and creating new records.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(client, args[0], shouldCreate, shouldDelete, shouldUpdate)
		},
	}

	deploy.Flags().BoolVarP(&shouldCreate, "create", "c", false, "create new records")
	deploy.Flags().BoolVarP(&shouldDelete, "delete", "d", false, "delete old records")
	deploy.Flags().BoolVarP(&shouldUpdate, "update", "u", false, "update changed records")

	return deploy
}

func deploy(client *porkbun.Client, configFile string, shouldCreate bool, shouldDelete bool, shouldUpdate bool) error {
	config, err := config.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("reading %v: %v", configFile, err)
	}

	from, err := client.AllRecords(config.Domain)
	if err != nil {
		return fmt.Errorf("fetching existing records: %v", err)
	}

	to := make([]porkbun.Record, len(config.Records))
	for i, record := range config.Records {
		to[i] = record.ToPorkbun()
	}

	added, removed, updated, unchanged := porkbun.DiffRecords(from, to)
	operations := buildOperations(added, removed, updated, unchanged, shouldCreate, shouldDelete, shouldUpdate)
	results := executeOperations(client, config.Domain, operations)
	printReport(results, shouldCreate, shouldDelete, shouldUpdate)

	for _, result := range results {
		if result.Status == deployment.Failure {
			return fmt.Errorf("couldn't %v record: %v", result.Type, result.Error)
		}
	}
	return nil
}

func buildOperations(added, removed, updated, unchanged []porkbun.Record, shouldCreate, shouldDelete, shouldUpdate bool) []deployment.RecordOperation {
	var operations []deployment.RecordOperation
	for _, record := range removed {
		operations = append(operations, deployment.RecordOperation{Type: deployment.Delete, Record: record, Planned: !shouldDelete})
	}
	for _, record := range updated {
		operations = append(operations, deployment.RecordOperation{Type: deployment.Update, Record: record, Planned: !shouldUpdate})
	}
	for _, record := range added {
		operations = append(operations, deployment.RecordOperation{Type: deployment.Create, Record: record, Planned: !shouldCreate})
	}
	for _, record := range unchanged {
		operations = append(operations, deployment.RecordOperation{Type: deployment.Keep, Record: record})
	}
	return operations
}

func executeOperations(client *porkbun.Client, domain string, operations []deployment.RecordOperation) []deployment.RecordOperationResult {
	var results []deployment.RecordOperationResult
	for _, operation := range operations {
		if operation.Planned {
			results = append(results, deployment.RecordOperationResult{
				Status: deployment.Planned,
				Type:   operation.Type,
				Record: operation.Record,
			})
			continue
		}

		result := operation.Execute(client, domain)
		results = append(results, result)
		if result.Status == deployment.Failure {
			break
		}
	}
	return results
}

func printReport(results []deployment.RecordOperationResult, shouldCreate, shouldDelete, shouldUpdate bool) {
	printSection(results, deployment.Delete, deleteSymbol, "Deleting", "Would delete", shouldDelete)
	printSection(results, deployment.Update, updateSymbol, "Updating", "Would update", shouldUpdate)
	printSection(results, deployment.Create, createSymbol, "Creating", "Would create", shouldCreate)

	fmt.Println("Keeping", countByType(results, deployment.Keep), "records:")
	for _, result := range results {
		if result.Type == deployment.Keep {
			fmt.Println(keepSymbol, result.Record)
		}
	}

	fullDeployment := shouldCreate && shouldDelete && shouldUpdate
	partialDeployment := shouldCreate || shouldDelete || shouldUpdate
	if fullDeployment {
		fmt.Println("Deployment complete!")
	} else if partialDeployment {
		fmt.Println("Partial deployment complete!")
	} else {
		fmt.Println("Mock deployment complete")
	}
}

func printSection(results []deployment.RecordOperationResult, operationType deployment.OperationType, symbol, executedVerb, plannedVerb string, executed bool) {
	count := countByType(results, operationType)
	if executed {
		fmt.Println(executedVerb, count, "records...")
	} else {
		fmt.Println(plannedVerb, count, "records:")
	}

	for _, result := range results {
		if result.Type != operationType {
			continue
		}
		if result.Status == deployment.Failure {
			fmt.Println(symbol, result.Record, "(failed:", fmt.Sprint(result.Error)+")")
		} else {
			fmt.Println(symbol, result.Record)
		}
	}
}

func countByType(results []deployment.RecordOperationResult, operationType deployment.OperationType) int {
	count := 0
	for _, result := range results {
		if result.Type == operationType {
			count++
		}
	}
	return count
}

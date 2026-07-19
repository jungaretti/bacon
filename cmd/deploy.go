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
	var dryRun bool
	var force bool

	deploy := &cobra.Command{
		Use:   "deploy <config-file>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by updating records in place when
possible, then deleting old records and creating new records. Previews the
deployment unless --force is specified.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(client, args[0], dryRun, force)
		},
	}

	deploy.Flags().BoolVarP(&dryRun, "dryrun", "d", false, "preview the deployment without making changes")
	deploy.Flags().BoolVarP(&force, "force", "f", false, "execute the deployment without confirmation")
	deploy.MarkFlagsMutuallyExclusive("dryrun", "force")

	return deploy
}

func deploy(client *porkbun.Client, configFile string, dryRun, force bool) error {
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
	recordDeployment := deployment.NewRecordDeployment(added, removed, updated, unchanged)

	var results []deployment.RecordOperationResult
	if force {
		results = recordDeployment.Execute(client, config.Domain)
	} else if dryRun {
		results = recordDeployment.Preview()
	} else {
		// Default to dry run
		results = recordDeployment.Preview()
	}
	printReport(results)

	for _, result := range results {
		if result.Status == deployment.Failure {
			return fmt.Errorf("couldn't %v record: %v", result.Type, result.Error)
		}
	}

	return nil
}

func printReport(results []deployment.RecordOperationResult) {
	executed := false
	for _, result := range results {
		if result.Status != deployment.Planned {
			executed = true
		}
	}

	printSection(results, deployment.Delete, deleteSymbol, "Deleting", "Would delete", executed)
	printSection(results, deployment.Update, updateSymbol, "Updating", "Would update", executed)
	printSection(results, deployment.Create, createSymbol, "Creating", "Would create", executed)

	fmt.Println("Keeping", countByType(results, deployment.Keep), "records:")
	for _, result := range results {
		if result.Type == deployment.Keep {
			fmt.Println(keepSymbol, result.Record)
		}
	}

	if executed {
		fmt.Println("Deployment complete!")
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

package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/deployment"
	"bacon/pkg/porkbun"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDeployCmd(client *porkbun.Client) *cobra.Command {
	var shouldCreate bool
	var shouldDelete bool
	var shouldUpdate bool
	var output string

	deploy := &cobra.Command{
		Use:   "deploy <config-file>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by updating records in place when
possible, then deleting old records and creating new records.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(client, args[0], shouldCreate, shouldDelete, shouldUpdate, output)
		},
	}

	deploy.Flags().BoolVarP(&shouldCreate, "create", "c", false, "create new records")
	deploy.Flags().BoolVarP(&shouldDelete, "delete", "d", false, "delete old records")
	deploy.Flags().BoolVarP(&shouldUpdate, "update", "u", false, "update changed records")
	deploy.Flags().StringVarP(&output, "output", "o", "table", "output format: table or json")

	return deploy
}

func deploy(client *porkbun.Client, configFile string, shouldCreate bool, shouldDelete bool, shouldUpdate bool, output string) error {
	renderer, err := deployment.NewDeploymentRenderer(output, os.Stdout)
	if err != nil {
		return err
	}

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

	operations := []deployment.Operation{}
	operations = appendOperations(operations, deployment.Delete, !shouldDelete, removed)
	operations = appendOperations(operations, deployment.Update, !shouldUpdate, updated)
	operations = appendOperations(operations, deployment.Create, !shouldCreate, added)
	operations = appendOperations(operations, deployment.Keep, false, unchanged)

	err = renderer.Preview(operations)
	if err != nil {
		return err
	}

	if shouldDelete {
		for _, record := range removed {
			err := client.DeleteRecord(config.Domain, record)
			if err != nil {
				return fmt.Errorf("couldn't delete record: %v", err)
			}
		}
	}

	if shouldUpdate {
		for _, record := range updated {
			err := client.EditRecord(config.Domain, record)
			if err != nil {
				return fmt.Errorf("couldn't update record: %v", err)
			}
		}
	}

	if shouldCreate {
		for _, record := range added {
			_, err := client.CreateRecord(config.Domain, record)
			if err != nil {
				return fmt.Errorf("couldn't create record: %v", err)
			}
		}
	}

	return renderer.Report(operations)
}

func appendOperations(results []deployment.Operation, action deployment.Action, dryRun bool, records []porkbun.Record) []deployment.Operation {
	for _, record := range records {
		results = append(results, deployment.Operation{Action: action, DryRun: dryRun, Record: record})
	}

	return results
}

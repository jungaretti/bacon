package cmd

import (
	"bacon/pkg/collections"
	"bacon/pkg/config"
	"bacon/pkg/operation"
	"bacon/pkg/porkbun"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newDeployCmd(client *porkbun.Client) *cobra.Command {
	var shouldCreate bool
	var shouldDelete bool
	var output string

	deploy := &cobra.Command{
		Use:   "deploy <config-file>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by deleting existing records and
creating new records.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reporter, err := newReporter(output, cmd.OutOrStdout())
			if err != nil {
				return err
			}

			return deploy(client, args[0], shouldCreate, shouldDelete, reporter)
		},
	}

	deploy.Flags().BoolVarP(&shouldCreate, "create", "c", false, "create new records")
	deploy.Flags().BoolVarP(&shouldDelete, "delete", "d", false, "delete old records")
	deploy.Flags().StringVarP(&output, "output", "o", "table", "output format: table or json")

	return deploy
}

func newReporter(output string, w io.Writer) (operation.Reporter, error) {
	switch output {
	case "table":
		return operation.NewTableReporter(w), nil
	case "json":
		return operation.NewJSONReporter(w), nil
	default:
		return nil, fmt.Errorf("unknown output format: %v", output)
	}
}

func deploy(client *porkbun.Client, configFile string, shouldCreate bool, shouldDelete bool, reporter operation.Reporter) error {
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

	added, removed, unchanged := collections.AddedRemovedUnchangedByHash(from, to, porkbun.RecordHash)

	var ops []operation.RecordOperation
	for _, record := range removed {
		ops = append(ops, operation.RecordOperation{
			Action:  operation.Delete,
			DryRun:  !shouldDelete,
			Record:  record,
			Execute: func() error { return client.DeleteRecord(config.Domain, record) },
		})
	}
	for _, record := range added {
		ops = append(ops, operation.RecordOperation{
			Action:  operation.Create,
			DryRun:  !shouldCreate,
			Record:  record,
			Execute: func() error { return client.CreateRecord(config.Domain, record) },
		})
	}
	for _, record := range unchanged {
		ops = append(ops, operation.RecordOperation{
			Action: operation.Skip,
			Record: record,
		})
	}

	manager := operation.NewManager()
	for _, op := range ops {
		reporter.Report(manager.Run(op))
	}

	summary := manager.Summary()
	err = reporter.Finish(summary)
	if err != nil {
		return fmt.Errorf("writing summary: %v", err)
	}

	if summary.Failed > 0 {
		return fmt.Errorf("%d operations failed", summary.Failed)
	}
	return nil
}

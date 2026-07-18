package cmd

import (
	"bacon/pkg/collections"
	"bacon/pkg/config"
	"bacon/pkg/porkbun"
	"fmt"

	"github.com/spf13/cobra"
)

func newDeployCmd(client *porkbun.Client) *cobra.Command {
	var shouldCreate bool
	var shouldDelete bool

	deploy := &cobra.Command{
		Use:   "deploy <config-file>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by deleting existing records and
creating new records.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(client, args[0], shouldCreate, shouldDelete)
		},
	}

	deploy.Flags().BoolVarP(&shouldCreate, "create", "c", false, "create new records")
	deploy.Flags().BoolVarP(&shouldDelete, "delete", "d", false, "delete old records")

	return deploy
}

func deploy(client *porkbun.Client, configFile string, shouldCreate bool, shouldDelete bool) error {
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

	added, removed := collections.AddedRemovedByHash(from, to, porkbun.RecordHash)
	if shouldDelete {
		fmt.Println("Deleting", len(removed), "records...")
		for _, record := range removed {
			err := client.DeleteRecord(config.Domain, record)
			if err != nil {
				return fmt.Errorf("couldn't delete record: %v", err)
			}
			fmt.Println("-", record)
		}
	} else {
		fmt.Println("Would delete", len(removed), "records:")
		for _, record := range removed {
			fmt.Println("-", record)
		}
	}
	if shouldCreate {
		fmt.Println("Creating", len(added), "records...")
		for _, record := range added {
			err := client.CreateRecord(config.Domain, record)
			if err != nil {
				return fmt.Errorf("couldn't create record: %v", err)
			}
			fmt.Println("-", record)
		}
	} else {
		fmt.Println("Would create", len(added), "records:")
		for _, record := range added {
			fmt.Println("-", record)
		}
	}

	if shouldCreate && shouldDelete {
		fmt.Println("Deployment complete!")
	} else if shouldCreate || shouldDelete {
		fmt.Println("Partial deployment complete!")
	} else {
		fmt.Println("Mock deployment complete")
	}
	return nil
}

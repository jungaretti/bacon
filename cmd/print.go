package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/porkbun"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newPrintCmd(client *porkbun.Client) *cobra.Command {
	print := &cobra.Command{
		Use:   "print <domain>",
		Short: "Print existing records",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return print(client, args[0])
		},
	}
	return print
}

func print(client *porkbun.Client, domain string) error {
	records, err := client.AllRecords(domain)
	if err != nil {
		return err
	}

	configRecords := make([]config.Record, len(records))
	for i, record := range records {
		configRecord, err := config.RecordFromPorkbun(record)
		if err != nil {
			return err
		}
		configRecords[i] = configRecord
	}

	config := config.Config{
		Domain:  domain,
		Records: configRecords,
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fmt.Printf("%s", out)
	return nil
}

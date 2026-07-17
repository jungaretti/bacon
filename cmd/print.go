package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/porkbun"
	"fmt"
	"strconv"

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
		ttl, err := strconv.Atoi(record.TTL)
		if err != nil {
			return fmt.Errorf("record %v has invalid TTL: %v", record, err)
		}

		configRecords[i] = config.Record{
			Name: record.Name,
			Type: record.Type,
			Ttl:  ttl,
			Data: record.Content,
		}
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

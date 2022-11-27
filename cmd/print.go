package cmd

import (
	"bacon/pkg/dns"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newPrintCmd(app *App) *cobra.Command {
	print := &cobra.Command{
		Use:   "print",
		Short: "Print existing records",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return print(app, args[0])
		},
	}
	return print
}

func print(app *App, domain string) error {
	records, err := app.Provider.AllRecords(domain)
	if err != nil {
		return err
	}

	configRecords := make([]dns.ConfigRecord, len(records))
	for i, record := range records {
		configRecords[i] = dns.ConfigFromRecord(record)
	}

	config := dns.Config{
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

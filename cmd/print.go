package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/dns"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newPrintCmd(provider dns.Provider) *cobra.Command {
	print := &cobra.Command{
		Use:   "print <domain>",
		Short: "Print existing records",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return print(provider, args[0])
		},
	}
	return print
}

func print(provider dns.Provider, domain string) error {
	records, err := provider.AllRecords(domain)
	if err != nil {
		return err
	}

	configRecords := make([]config.Record, len(records))
	for i, record := range records {
		// Sets ttl to nil if Atoi fails
		ttl, _ := strconv.Atoi(record.GetTtl())

		configRecords[i] = config.Record{
			Name: record.GetName(),
			Type: record.GetType(),
			Ttl:  &ttl,
			Data: record.GetData(),
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

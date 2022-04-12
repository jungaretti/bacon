package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDumpCmd(app *App) *cobra.Command {
	dump := &cobra.Command{
		Use:   "dump <domain>",
		Short: "Fetch all of a domain's DNS records",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dump(app, args[0])
		},
	}
	return dump
}

func dump(app *App, domain string) error {
	records, err := app.Client.GetRecords(domain)
	if err != nil {
		return err
	}

	fmt.Println(records)
	return nil
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDumpCmd(app *App) *cobra.Command {
	dump := &cobra.Command{
		Use:   "dump <domain>",
		Short: "Retrieve DNS records for domain",
		Long:  `Retrieve all editable DNS records associated with a domain.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dump(app, args[0])
		},
	}
	return dump
}

func dump(app *App, domain string) {
	msg, err := app.Client.GetRecords(domain)
	if err != nil {
		errMsg := fmt.Errorf("error sending request: %w", err)
		fmt.Println(errMsg)
		return
	}

	fmt.Println(msg)
}

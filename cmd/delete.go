package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDeleteCmd(app *App) *cobra.Command {
	delete := &cobra.Command{
		Use:   "delete <domain> <id>",
		Short: "Delete existing DNS record for domain",
		Long:  `Delete an existing DNS record for a domain.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			delete(app, args[0], args[1])
		},
	}
	return delete
}

func delete(app *App, domain string, id string) {
	msg, err := app.Client.DeleteRecord(domain, id)
	if err != nil {
		errMsg := fmt.Errorf("error deleting record: %w", err)
		fmt.Println(errMsg)
		return
	}

	fmt.Println(msg)
}

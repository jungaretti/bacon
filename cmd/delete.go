package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDeleteCmd(app *App) *cobra.Command {
	delete := &cobra.Command{
		Use:   "delete <domain> <id>",
		Short: "Delete an existing DNS record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(app, args[0], args[1])
		},
	}
	return delete
}

func delete(app *App, domain string, id string) error {
	ack, err := app.Client.DeleteRecord(domain, id)
	if err != nil {
		return err
	}

	if ack.Ok {
		fmt.Println(ack.Message)
		return nil
	} else {
		return fmt.Errorf(ack.Message)
	}
}

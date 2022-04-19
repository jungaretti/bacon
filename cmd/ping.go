package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPingCmd(app *App) *cobra.Command {
	ping := &cobra.Command{
		Use:   "ping",
		Short: "Say hello to Porkbun",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ping(app)
		},
	}
	return ping
}

func ping(app *App) error {
	ack, err := app.Client.Ping()
	if err != nil {
		return err
	}

	if ack.Ok {
		fmt.Printf("Success! %s is ready to use.\n", app.Client.GetName())
		return nil
	} else {
		return fmt.Errorf(ack.Message)
	}
}

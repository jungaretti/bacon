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
		fmt.Println("Success! %s is ready to use.", app.Client.Name())
		return nil
	} else {
		return fmt.Errorf(ack.Message)
	}
}

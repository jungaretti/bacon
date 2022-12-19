package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPingCmd(app *App) *cobra.Command {
	ping := &cobra.Command{
		Use:   "ping",
		Short: "Check authentication status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Ping(app)
		},
	}
	return ping
}

func Ping(app *App) error {
	err := app.Provider.CheckAuth()
	if err != nil {
		return err
	}

	fmt.Printf("Success! %s is ready to use.\n", "Porkbun")
	return nil
}

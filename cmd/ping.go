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
		Run: func(cmd *cobra.Command, args []string) {
			ping(app)
		},
	}
	return ping
}

func ping(app *App) {
	msg, err := app.Client.Ping()
	if err != nil {
		errMsg := fmt.Errorf("error sending request: %w", err)
		fmt.Println(errMsg)
		return
	}

	fmt.Println(msg)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPingCmd(app *App) *cobra.Command {
	ping := &cobra.Command{
		Use:   "ping",
		Short: "Say hello to Porkbun",
		Long:  `You can test communication with the API using the ping endpoint.`,
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

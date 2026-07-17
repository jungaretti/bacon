package cmd

import (
	"bacon/pkg/porkbun"
	"fmt"

	"github.com/spf13/cobra"
)

func newPingCmd(client *porkbun.Client) *cobra.Command {
	ping := &cobra.Command{
		Use:   "ping",
		Short: "Check authentication status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ping(client)
		},
	}
	return ping
}

func ping(client *porkbun.Client) error {
	err := client.Ping()
	if err != nil {
		return err
	}

	fmt.Printf("Success! %s is ready to use.\n", "Porkbun")
	return nil
}

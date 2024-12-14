package cmd

import (
	"bacon/pkg/dns"
	"fmt"

	"github.com/spf13/cobra"
)

func newPingCmd(provider dns.Provider) *cobra.Command {
	ping := &cobra.Command{
		Use:   "ping",
		Short: "Check authentication status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ping(provider)
		},
	}
	return ping
}

func ping(provider dns.Provider) error {
	err := provider.CheckAuth()
	if err != nil {
		return err
	}

	fmt.Printf("Success! %s is ready to use.\n", "Porkbun")
	return nil
}

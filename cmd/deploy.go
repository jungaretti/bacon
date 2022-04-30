package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDeployCmd(app *App) *cobra.Command {
	var create bool
	var delete bool

	deploy := &cobra.Command{
		Use:   "deploy <domain> <config>",
		Short: "Deploy DNS records from a config file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(app, args[0], args[1], &create, &delete)
		},
	}

	deploy.Flags().BoolVarP(&create, "create", "c", false, "create new records")
	deploy.Flags().BoolVarP(&delete, "delete", "d", false, "delete old records")

	return deploy
}

func deploy(app *App, domain string, configFile string, create *bool, delete *bool) error {
	return fmt.Errorf("haven't implemented deploy yet")
}

package cmd

import (
	"bacon/client"
	"fmt"

	"github.com/spf13/cobra"
)

func newDeployCmd(app *App) *cobra.Command {
	var create bool
	var delete bool

	deploy := &cobra.Command{
		Use:   "deploy <domain> <config>",
		Short: "Deploy DNS records from a file",
		Long:  `Deploy all DNS records from a config file.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			deploy(app, args[0], args[1], &create, &delete)
		},
	}

	deploy.Flags().BoolVarP(&create, "create", "c", false, "create new records")
	deploy.Flags().BoolVarP(&delete, "delete", "d", false, "delete old records")

	return deploy
}

func deploy(app *App, domain string, configFile string, create *bool, delete *bool) {
	config, err := client.ReadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !*create {
		fmt.Println("Not creating new records!")
	}
	if !*delete {
		fmt.Println("Not deleting old records!")
	}

	ack, err := app.Client.SyncRecords(domain, config.Records, *create, *delete)
	if err != nil {
		errMsg := fmt.Errorf("error sending request: %w", err)
		fmt.Println(errMsg)
		return
	}

	fmt.Println(ack)
}

package cmd

import (
	"bacon/client"
	"fmt"

	"github.com/spf13/cobra"
)

func newDeployCmd(app *App) *cobra.Command {
	deploy := &cobra.Command{
		Use:   "deploy <domain> <config>",
		Short: "Deploy DNS records from a file",
		Long:  `Deploy all DNS records from a config file.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			deploy(app, args[0], args[1])
		},
	}
	return deploy
}

func deploy(app *App, domain string, configFile string) {
	config, err := client.ReadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Delete old records
	// TODO: Create new records

	for _, record := range config.Records {
		msg, err := app.Client.CreateRecord(domain, &record)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(msg)
	}
}

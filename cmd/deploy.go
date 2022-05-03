package cmd

import (
	"bacon/pkg/client"
	"bacon/pkg/helpers"
	"fmt"

	"github.com/spf13/cobra"
)

func newDeployCmd(app *App) *cobra.Command {
	var shouldCreate bool
	var shouldDelete bool

	deploy := &cobra.Command{
		Use:   "deploy <config>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by deleting existing records and
creating new records.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(app, args[0], shouldCreate, shouldDelete)
		},
	}

	deploy.Flags().BoolVarP(&shouldCreate, "create", "c", false, "create new records")
	deploy.Flags().BoolVarP(&shouldDelete, "delete", "d", false, "delete old records")

	return deploy
}

func deploy(app *App, configFile string, shouldCreate bool, shouldDelete bool) error {
	config := client.Config{}
	err := helpers.ReadAndParseYamlFile(configFile, &config)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return app.Client.Deploy(config.Domain, config.Records, shouldCreate, shouldDelete)
}

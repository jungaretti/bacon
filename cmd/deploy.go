package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/deployment"
	"bacon/pkg/porkbun"
	"fmt"

	"github.com/spf13/cobra"
)

func newDeployCmd(client *porkbun.Client) *cobra.Command {
	var dryRun bool
	var force bool

	deploy := &cobra.Command{
		Use:   "deploy <config-file>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by updating records in place when
possible, then deleting old records and creating new records. Previews the
deployment unless --force is specified.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(client, args[0], force)
		},
	}

	deploy.Flags().BoolVar(&dryRun, "dry-run", true, "preview the deployment without making changes")
	deploy.Flags().BoolVar(&force, "force", false, "execute the deployment without confirmation")
	deploy.MarkFlagsMutuallyExclusive("dry-run", "force")

	return deploy
}

func deploy(client *porkbun.Client, configFile string, force bool) error {
	config, err := config.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("reading %v: %v", configFile, err)
	}

	from, err := client.AllRecords(config.Domain)
	if err != nil {
		return fmt.Errorf("fetching existing records: %v", err)
	}

	to := make([]porkbun.Record, len(config.Records))
	for i, record := range config.Records {
		to[i] = record.ToPorkbun()
	}

	added, removed, updated, unchanged := porkbun.DiffRecords(from, to)
	recordDeployment := deployment.NewDeployment(added, removed, updated, unchanged)

	var deploymentResult deployment.DeploymentResult
	if force {
		deploymentResult = recordDeployment.Execute(client, config.Domain)
	} else {
		deploymentResult = recordDeployment.Preview()
	}

	formatter := deployment.TextFormatter{}
	fmt.Print(formatter.Format(deploymentResult))

	for _, result := range deploymentResult.Results {
		if result.Status == deployment.Failure {
			return fmt.Errorf("couldn't %v record: %v", result.Type, result.Error)
		}
	}

	return nil
}

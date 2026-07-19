package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/deployment"
	"bacon/pkg/porkbun"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newDeployCmd(client *porkbun.Client) *cobra.Command {
	var dryRun bool
	var force bool
	var confirm bool
	var output string

	deploy := &cobra.Command{
		Use:   "deploy <config-file>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by updating records in place when
possible, then deleting old records and creating new records. Previews the
deployment unless --force is specified.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(client, args[0], force, confirm, output)
		},
	}

	deploy.Flags().BoolVar(&dryRun, "dry-run", true, "preview the deployment without making changes")
	deploy.Flags().BoolVar(&force, "force", false, "execute the deployment without confirmation")
	deploy.Flags().BoolVar(&confirm, "confirm", false, "preview and confirm the deployment")
	deploy.MarkFlagsMutuallyExclusive("dry-run", "force", "confirm")

	deploy.Flags().StringVarP(&output, "output", "o", "table", "output format: table or json")

	return deploy
}

func deploy(client *porkbun.Client, configFile string, force bool, confirm bool, output string) error {
	var formatter deployment.Formatter
	switch output {
	case "table":
		formatter = deployment.TableFormatter{}
	case "json":
		formatter = deployment.JSONFormatter{}
	default:
		return fmt.Errorf("unknown output format: %v", output)
	}

	if confirm && output == "json" {
		return errors.New("--confirm is not compatible with --output json")
	}

	config, err := config.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("reading %v: %v", configFile, err)
	}

	fmt.Print(formatter.FormatStart(config.Domain, !force))

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

	switch {
	case confirm:
		previewResult := recordDeployment.Preview()
		fmt.Print(formatter.FormatResult(previewResult))

		counts := previewResult.Summary.OperationCounts
		if counts[deployment.Delete]+counts[deployment.Update]+counts[deployment.Create] == 0 {
			fmt.Println("No changes to deploy")
			return nil
		}

		if !promptProceed(os.Stdin) {
			fmt.Println("Deployment cancelled")
			return nil
		}

		fmt.Print(formatter.FormatStart(config.Domain, false))
		deploymentResult := recordDeployment.Execute(client, config.Domain)
		fmt.Print(formatter.FormatResult(deploymentResult))

		return firstFailure(deploymentResult)
	case force:
		deploymentResult := recordDeployment.Execute(client, config.Domain)
		fmt.Print(formatter.FormatResult(deploymentResult))

		return firstFailure(deploymentResult)
	default:
		fmt.Print(formatter.FormatResult(recordDeployment.Preview()))
		return nil
	}
}

func promptProceed(reader io.Reader) bool {
	fmt.Print("Proceed with deployment? [y/N]: ")

	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		return false
	}

	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return answer == "y" || answer == "yes"
}

func firstFailure(deploymentResult deployment.DeploymentResult) error {
	for _, result := range deploymentResult.Results {
		if result.Status == deployment.Failure {
			return errors.New(result.Error)
		}
	}
	return nil
}

package cmd

import (
	"bacon/pkg/dns"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	config, err := readConfig(configFile)
	if err != nil {
		return fmt.Errorf("reading config: %v", err)
	}

	from, err := app.Provider.AllRecords(config.Domain)
	if err != nil {
		return fmt.Errorf("fetching existing records: %v", err)
	}

	configRecords := config.Records
	to := make([]dns.Record, len(configRecords))
	for i, record := range configRecords {
		to[i] = record
	}

	added, removed := dns.DiffRecords(from, to)
	if shouldDelete {
	} else {
		fmt.Println("Would delete", len(removed), "records:")
		for _, record := range removed {
			fmt.Println("-", record)
		}
	}
	if shouldCreate {
	} else {
		fmt.Println("Would create", len(added), "records:")
		for _, record := range added {
			fmt.Println("-", record)
		}
	}

	return nil
}

func readConfig(configFile string) (*dns.Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	config := dns.Config{}
	err = yaml.Unmarshal(raw, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

package cmd

import (
	"bacon/client"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy <domain> <config>",
	Short: "Deploy DNS records from a file",
	Long:  `Deploy all DNS records from a config file.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		deploy(args[0], args[1])
	},
}

type Config struct {
	Records []client.Record `yaml:"records"`
}

func deploy(domain string, configFile string) {
	pork := client.Pork{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	config, err := client.ReadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Delete old records
	// TODO: Create new records

	for _, record := range config.Records {
		msg, err := pork.CreateRecord(domain, &record)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(msg)
	}
}

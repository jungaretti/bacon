package cmd

import (
	"bacon/client"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var record = client.Record{}

func init() {
	createCmd.Flags().StringVarP(&record.Type, "type", "t", "", "type of DNS record")
	createCmd.Flags().StringVarP(&record.Host, "host", "H", "", "subdomain of DNS record")
	createCmd.Flags().StringVarP(&record.Content, "content", "c", "", "content of DNS record")
	createCmd.Flags().IntVarP(&record.TTL, "ttl", "l", 300, "TTL of DNS record")
	createCmd.Flags().IntVarP(&record.Priority, "priority", "p", 20, "priority of DNS record")

	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("content")

	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create <domain>",
	Short: "Create new DNS record for domain",
	Long:  `Create a new DNS record for a domain.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		create(args[0])
	},
}

func create(domain string) {
	pork := client.Pork{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	msg, err := pork.CreateRecord(domain, &record)
	if err != nil {
		errMsg := fmt.Errorf("error creating record: %w", err)
		fmt.Println(errMsg)
	}

	fmt.Println(msg)
}

package cmd

import (
	"bacon/client"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump <domain>",
	Short: "Retrieve DNS records for a domain",
	Long:  `Retrieve all editable DNS records associated with a domain.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dump(args[0])
	},
}

func dump(domain string) {
	auth := client.Auth{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	records, err := client.GetRecords(&auth, domain)
	if err != nil {
		fmt.Println(fmt.Errorf("couldn't dump from Porkbun: %w", err))
	}

	for _, element := range *records {
		fmt.Printf("%s: Name=%s, Type=%s, Content=%s, TTL=%s, Priority=%s, Notes=%s\n", element.Id, element.Name, element.Type, element.Content, element.Ttl, element.Priority, element.Notes)
	}
}

package cmd

import (
	"bacon/client"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete <domain> <id>",
	Short: "Delete existing DNS record for domain",
	Long:  `Delete an existing DNS record for a domain.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		delete(args[0], args[1])
	},
}

func delete(domain string, id string) {
	auth := client.Auth{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	msg, err := client.DeleteRecordJSON(&auth, domain, id)
	if err != nil {
		errMsg := fmt.Errorf("error deleting record: %w", err)
		fmt.Println(errMsg)
	}

	fmt.Printf("%s\n", msg)
}

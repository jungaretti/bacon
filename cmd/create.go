package cmd

import (
	"bacon/client"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCreateCmd(app *App) *cobra.Command {
	record := client.Record{}

	create := &cobra.Command{
		Use:   "create <domain>",
		Short: "Create a new DNS record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			create(app, args[0], &record)
		},
	}

	create.Flags().StringVarP(&record.Type, "type", "t", "", "type of DNS record")
	create.Flags().StringVarP(&record.Host, "host", "H", "", "subdomain of DNS record")
	create.Flags().StringVarP(&record.Content, "content", "c", "", "content of DNS record")
	create.Flags().IntVarP(&record.TTL, "ttl", "l", 300, "TTL of DNS record")
	create.Flags().IntVarP(&record.Priority, "priority", "p", 20, "priority of DNS record")

	create.MarkFlagRequired("type")
	create.MarkFlagRequired("content")

	return create
}

func create(app *App, domain string, record *client.Record) {
	pork := client.Pork{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	msg, err := pork.CreateRecord(domain, record)
	if err != nil {
		errMsg := fmt.Errorf("error creating record: %w", err)
		fmt.Println(errMsg)
		return
	}

	fmt.Println(msg)
}

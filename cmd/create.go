package cmd

import (
	"bacon/client"
	"bacon/client/porkbun"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(app, args[0], &record)
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

func create(app *App, domain string, record *client.Record) error {
	pork := porkbun.PorkClient{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	ack, err := pork.CreateRecord(domain, record)
	if err != nil {
		return err
	}

	if ack.Ok {
		fmt.Println(ack.Message)
		return nil
	} else {
		return fmt.Errorf(ack.Message)
	}
}

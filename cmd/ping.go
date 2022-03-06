package cmd

import (
	"bacon/client"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Say hello to Porkbun",
	Long:  `You can test communication with the API using the ping endpoint.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ping()
	},
}

func ping() {
	auth := client.Auth{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	ack, err := client.Ping(&auth)
	if err != nil {
		msg := fmt.Errorf("couldn't ping Porkbun: %w", err)
		fmt.Println(msg)
	}

	if ack.Success {
		fmt.Println("Success!")
	} else {
		fmt.Println("Failure.")
	}
}

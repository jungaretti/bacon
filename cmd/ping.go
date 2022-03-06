package cmd

import (
	"bacon/pork"
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
	Long: `You can test communication with the API using the ping endpoint.
The ping endpoint will also return your IP address,
this can be handy when building dynamic DNS clients.`,
	Run: func(cmd *cobra.Command, args []string) {
		ping()
	},
}

func ping() {
	auth := pork.Auth{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	msg, err := pork.Ping(auth)
	if err != nil {
		msg := fmt.Errorf("couldn't ping Porkbun: %w", err)
		fmt.Println(msg)
	}

	fmt.Println(msg)
}

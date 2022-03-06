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
	Long:  `Does this really need a long description?`,
	Run: func(cmd *cobra.Command, args []string) {
		executePing()
	},
}

func executePing() {
	auth := pork.Auth{
		ApiKey:       os.Getenv("PORKBUN_API_KEY"),
		SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
	}

	msg, err := pork.Ping(auth)
	if err != nil {
		fmt.Errorf("couldn't ping Porkbun: %w", err)
	}

	fmt.Println(msg)
}

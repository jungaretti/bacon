package cmd

import (
	"bacon/pork"
	"fmt"

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
		ping()
	},
}

func ping() {
	msg, err := pork.Ping()
	if err != nil {
		fmt.Errorf("couldn't ping Porkbun: %w", err)
	}

	fmt.Println(msg)
}

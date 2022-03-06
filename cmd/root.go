package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bacon",
	Short: "Bacon is a DNS manager",
	Long: `A flexible DNS record manager for Porkbun. Complete documentation
is available at https://github.com/jungaretti/bacon.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

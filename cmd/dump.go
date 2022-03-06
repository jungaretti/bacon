package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "dump",
	Short: "Export all DNS records",
	Long:  `Does this really need a long description?`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dumping DNS records...")
	},
}

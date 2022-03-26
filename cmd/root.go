package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute(app *App) {
	root := newRootCmd(app)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd(app *App) *cobra.Command {
	root := &cobra.Command{
		Use:   "bacon",
		Short: "Bacon is a fast and flexible DNS manager",
	}

	// Use command constructors to share one client
	root.AddCommand(newCreateCmd(app))
	root.AddCommand(newDeleteCmd(app))
	root.AddCommand(newDeployCmd(app))
	root.AddCommand(newDumpCmd(app))
	root.AddCommand(newPingCmd(app))

	return root
}

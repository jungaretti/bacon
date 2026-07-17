package cmd

import (
	"bacon/pkg/porkbun"
	"os"

	"github.com/spf13/cobra"
	"github.com/subosito/gotenv"
)

func Execute() {
	// Loads .env in the current directory
	gotenv.Load()

	var porkbunApiKey = os.Getenv("PORKBUN_API_KEY")
	var porkbunSecretApiKey = os.Getenv("PORKBUN_SECRET_KEY")
	var client = porkbun.NewClient(porkbunApiKey, porkbunSecretApiKey)

	root := newRootCmd(client)
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func newRootCmd(client *porkbun.Client) *cobra.Command {
	root := &cobra.Command{
		Use:   "bacon",
		Short: "Bacon is a tasty DNS manager for Porkbun",
	}

	root.AddCommand(newPingCmd(client))
	root.AddCommand(newDeployCmd(client))
	root.AddCommand(newPrintCmd(client))

	return root
}

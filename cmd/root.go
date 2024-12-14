package cmd

import (
	"bacon/pkg/dns"
	"bacon/pkg/providers/porkbun"
	"os"

	"github.com/spf13/cobra"
	"github.com/subosito/gotenv"
)

func Execute() {
	// Loads .env in the current directory
	gotenv.Load()

	var porkbunApiKey = os.Getenv("PORKBUN_API_KEY")
	var porkbunSecretApiKey = os.Getenv("PORKBUN_SECRET_KEY")
	var porkbunProvider = porkbun.NewPorkbunProvider(porkbunApiKey, porkbunSecretApiKey)

	root := newRootCmd(porkbunProvider)
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func newRootCmd(provider dns.Provider) *cobra.Command {
	root := &cobra.Command{
		Use:   "bacon",
		Short: "Bacon is a tasty DNS manager for Porkbun",
	}

	root.AddCommand(newPingCmd(provider))
	root.AddCommand(newDeployCmd(provider))
	root.AddCommand(newPrintCmd(provider))

	return root
}

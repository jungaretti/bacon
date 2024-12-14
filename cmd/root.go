package cmd

import (
	"bacon/pkg/providers/porkbun"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/subosito/gotenv"
)

func Execute() {
	app := newPorkbunApp()

	root := newRootCmd(app)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd(app *App) *cobra.Command {
	root := &cobra.Command{
		Use:   "bacon",
		Short: "Bacon is a tasty DNS manager for Porkbun",
	}

	// Use command constructors to share one app
	root.AddCommand(newPingCmd(app))
	root.AddCommand(newDeployCmd(app))
	root.AddCommand(newPrintCmd(app))

	return root
}

func newPorkbunApp() *App {
	// Loads .env in the current directory
	gotenv.Load()

	var porkbunApiKey = os.Getenv("PORKBUN_API_KEY")
	var porkbunSecretApiKey = os.Getenv("PORKBUN_SECRET_KEY")

	app := App{
		Provider: porkbun.NewPorkbunProvider(porkbunApiKey, porkbunSecretApiKey),
	}
	return &app
}

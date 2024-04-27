package main

import (
	"bacon/cmd"
	"bacon/pkg/providers/porkbun"
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	// Loads .env in the current directory
	gotenv.Load()

	var porkbunApiKey = os.Getenv("PORKBUN_API_KEY")
	var porkbunSecretApiKey = os.Getenv("PORKBUN_SECRET_KEY")

	app := cmd.App{
		Provider: porkbun.NewPorkbunProvider(porkbunApiKey, porkbunSecretApiKey),
	}

	cmd.Execute(&app)
}

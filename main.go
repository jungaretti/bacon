package main

import (
	"bacon/cmd"
	"bacon/pkg/providers/porkbun"
	"bacon/pkg/secrets"

	"github.com/subosito/gotenv"
)

func main() {
	// Loads .env in the current directory
	gotenv.Load()

	// Only supports Porkbun... for now :D
	app := cmd.App{
		Provider: porkbun.NewPorkbunProvider(secrets.Provider{}),
	}

	cmd.Execute(&app)
}

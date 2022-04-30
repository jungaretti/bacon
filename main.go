package main

import (
	"bacon/cmd"
	"bacon/pkg/porkbun"
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	// Loads .env in the current directory
	gotenv.Load()

	// Only supports Porkbun... for now :D
	app := cmd.App{
		Client: &porkbun.PorkClient{
			Auth: porkbun.PorkAuth{
				ApiKey:       os.Getenv("PORKBUN_API_KEY"),
				SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
			},
		},
	}

	cmd.Execute(&app)
}

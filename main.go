package main

import (
	"bacon/client"
	"bacon/cmd"
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	app := cmd.App{
		Client: client.Pork{
			ApiKey:       os.Getenv("PORKBUN_API_KEY"),
			SecretApiKey: os.Getenv("PORKBUN_SECRET_KEY"),
		},
	}

	cmd.Execute(&app)
}

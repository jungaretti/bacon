package porkbun

import "fmt"

type PorkClient struct {
	ApiKey       string
	SecretApiKey string
}

func (client *PorkClient) Name() string {
	return "Porkbun"
}

func (client *PorkClient) Ping() error {
	return fmt.Errorf("haven't implemented ping yet")
}

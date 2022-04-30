package porkbun

import "bacon/pkg/client"

type PorkAuth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

type PorkClient struct {
	Auth PorkAuth
}

func (client *PorkClient) Name() string {
	return "Porkbun"
}

func (client *PorkClient) Ping() error {
	return ping(client.Auth)
}

func (client *PorkClient) Deploy(domain string, records []client.Record, shouldCreate bool, shouldDelete bool) error {
	return deploy(client.Auth, domain, shouldCreate, shouldDelete)
}

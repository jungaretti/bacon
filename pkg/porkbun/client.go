package porkbun

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

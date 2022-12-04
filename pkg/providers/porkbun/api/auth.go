package api

type Auth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

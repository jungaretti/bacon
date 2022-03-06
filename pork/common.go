package pork

var PORKBUN_PING = "https://porkbun.com/api/json/v3/ping"

type Auth struct {
	ApiKey       string `json:"apikey"`
	SecretApiKey string `json:"secretapikey"`
}

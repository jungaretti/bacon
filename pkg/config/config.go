package config

type Config struct {
	Domain  string   `yaml:"domain"`
	Records []Record `yaml:"records"`
}

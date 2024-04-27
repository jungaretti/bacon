package config

type Config struct {
	Domain  string         `yaml:"domain"`
	Records []ConfigRecord `yaml:"records"`
}

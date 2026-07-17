package config

type Record struct {
	Name     string `yaml:"host"`
	Type     string `yaml:"type"`
	Ttl      int    `yaml:"ttl"`
	Data     string `yaml:"content"`
	Priority int    `yaml:"priority"`
}

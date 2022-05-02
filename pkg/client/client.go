package client

type Record struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Content  string `yaml:"content"`
	TTL      int    `yaml:"ttl"`
	Priority int    `yaml:"priority"`
	Notes    string `yaml:"notes"`
}

type Config struct {
	Domain  string   `yaml:"domain"`
	Records []Record `yaml:"records"`
}

type Client interface {
	Name() string
	Ping() error
	Deploy(domain string, records []Record, shouldCreate bool, shouldDelete bool) error
}

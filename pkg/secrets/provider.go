package secrets

import "os"

type Provider struct {
	local map[string]string
}

func (p *Provider) Set(key, value string) {
	if p.local == nil {
		p.local = make(map[string]string)
	}

	p.local[key] = value
}

func (p *Provider) Read(key string) string {
	if p.local == nil {
		p.local = make(map[string]string)
	}

	local := p.local[key]
	if local != "" {
		return local
	}

	return os.Getenv(key)
}

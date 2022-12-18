package mock

import "bacon/pkg/dns"

type MockRecord struct {
	Name string
	Type string
	Ttl  string
	Data string
}

func (r MockRecord) GetName() string {
	return r.Name
}

func (r MockRecord) GetType() string {
	return r.Type
}

func (r MockRecord) GetTtl() string {
	return r.Ttl
}

func (r MockRecord) GetData() string {
	return r.Data
}

var _ dns.Record = MockRecord{}

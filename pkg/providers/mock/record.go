package mock

import "bacon/pkg/dns"

type MockRecord struct {
	Name string
	Type string
	Ttl  string
	Data string
}

func (r MockRecord) GetId() string {
	return r.Name + r.Type
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

func NewMockRecord(record dns.Record) MockRecord {
	new := MockRecord{
		Name: record.GetName(),
		Type: record.GetType(),
		Ttl:  record.GetTtl(),
		Data: record.GetData(),
	}

	return new
}

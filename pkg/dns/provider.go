package dns

type Provider interface {
	CheckAuth() error
	AllRecords(domain string) ([]Record, error)
	CreateRecord(domain string, record Record) error
	DeleteRecord(domain string, record Record) error
}

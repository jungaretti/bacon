package dns

type Provider interface {
	CheckAuth() error
	AllRecords() ([]Record, error)
	CreateRecord(Record) error
	DeleteRecord(Record) error
}

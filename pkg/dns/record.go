package dns

type Record interface {
	GetType() string
	GetHost() string
	GetContent() string
	GetTTL() string
}

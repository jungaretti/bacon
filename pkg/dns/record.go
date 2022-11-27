package dns

type Record interface {
	GetType() string
	GetHost() string
	GetContent() string
	GetTTL() string
}

func RecordEquals(left Record, right Record) bool {
	equal := left.GetType() == right.GetType()
	equal = equal && left.GetHost() == right.GetHost()
	equal = equal && left.GetContent() == right.GetContent()
	equal = equal && left.GetTTL() == right.GetTTL()

	return equal
}

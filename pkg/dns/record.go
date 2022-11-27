package dns

// DNS record according to the RFC 1035 spec.
//
// [RFC 1035]: https://www.rfc-editor.org/rfc/rfc1035
type Record interface {
	GetName() string
	GetType() string
	GetTtl() string
	GetData() string
}

func RecordEquals(l Record, r Record) (equal bool) {
	equal = equal && l.GetName() == r.GetName()
	equal = equal && l.GetType() == r.GetType()
	equal = equal && l.GetTtl() == r.GetTtl()
	equal = equal && l.GetData() == r.GetData()
	return
}

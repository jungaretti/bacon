package dns

// DNS record according to the RFC 1035 spec.
//
// [RFC 1035]: https://www.rfc-editor.org/rfc/rfc1035
type Record interface {
	GetName() string
	GetType() string
	GetTtl() string
	GetData() string
	Equals(Record) bool
	Hash() string
}

func RecordEquals(l Record, r Record) bool {
	equal := true
	equal = equal && l.GetName() == r.GetName()
	equal = equal && l.GetType() == r.GetType()
	equal = equal && l.GetTtl() == r.GetTtl()
	equal = equal && l.GetData() == r.GetData()
	return equal
}

func RecordHash(r Record) string {
	return r.GetName() + r.GetType() + r.GetTtl() + r.GetData()
}

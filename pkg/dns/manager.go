package dns

func difference(from, to []Record) []Record {
	set := make(map[string]bool)

	for _, record := range to {
		set[record.Hash()] = true
	}

	var missing []Record
	for _, record := range from {
		if _, ok := set[record.Hash()]; !ok {
			missing = append(missing, record)
		}
	}
	return missing
}

func DiffRecords(from, to []Record) ([]Record, []Record) {
	added := difference(to, from)
	removed := difference(from, to)
	return added, removed
}

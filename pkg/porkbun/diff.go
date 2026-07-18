package porkbun

import "bacon/pkg/collections"

// Diff existing records (from) against desired records (to).
func DiffRecords(from, to []Record) (added, removed, edited, unchanged []Record) {
	added, removed, unchanged = collections.AddedRemovedUnchangedByHash(from, to, RecordHash)
	groups, removed, added := collections.GroupByHash(removed, added, RecordIdentityHash)

	edited = make([]Record, len(groups))
	for i, group := range groups {
		record := group.To
		record.Id = group.From.Id
		edited[i] = record
	}

	return added, removed, edited, unchanged
}

package porkbun

import "bacon/pkg/collections"

// Diff existing records (from) against desired records (to).
func DiffRecords(from, to []Record) (added, removed, updated, unchanged []Record) {
	added, removed, unchanged = collections.AddedRemovedUnchangedByHash(from, to, RecordHash)
	groups, removed, added := collections.GroupByHash(removed, added, RecordIdentityHash)

	updated = make([]Record, len(groups))
	for i, group := range groups {
		record := group.To
		record.Id = group.From.Id
		updated[i] = record
	}

	return added, removed, updated, unchanged
}

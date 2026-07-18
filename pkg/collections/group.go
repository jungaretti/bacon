package collections

type Group[T any] struct {
	From T
	To   T
}

// Group elements that share the same hash.
func GroupByHash[T any, U comparable](from, to []T, hasher func(T) U) (groups []Group[T], ungroupedFrom, ungroupedTo []T) {
	fromByHash := make(map[U][]T)
	for _, element := range from {
		hash := hasher(element)
		fromByHash[hash] = append(fromByHash[hash], element)
	}

	grouped := make(map[U]int)
	for _, element := range to {
		hash := hasher(element)
		if remaining := fromByHash[hash]; len(remaining) > 0 {
			groups = append(groups, Group[T]{From: remaining[0], To: element})
			fromByHash[hash] = remaining[1:]
			grouped[hash]++
		} else {
			ungroupedTo = append(ungroupedTo, element)
		}
	}

	skipped := make(map[U]int)
	for _, element := range from {
		hash := hasher(element)
		if skipped[hash] < grouped[hash] {
			skipped[hash]++
			continue
		}
		ungroupedFrom = append(ungroupedFrom, element)
	}

	return groups, ungroupedFrom, ungroupedTo
}

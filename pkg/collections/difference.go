package collections

func SetDifferenceByHash[T any, U comparable](from, to []T, hasher func(T) U) []T {
	set := make(map[U]bool)

	for _, element := range to {
		set[hasher(element)] = true
	}

	var missing []T
	for _, element := range from {
		if !set[hasher(element)] {
			missing = append(missing, element)
		}
	}
	return missing
}

func SetIntersectionByHash[T any, U comparable](from, to []T, hasher func(T) U) []T {
	set := make(map[U]bool)

	for _, element := range to {
		set[hasher(element)] = true
	}

	var intersection []T
	for _, element := range from {
		if set[hasher(element)] {
			intersection = append(intersection, element)
		}
	}
	return intersection
}

func AddedRemovedUnchangedByHash[T any, U comparable](from, to []T, hasher func(T) U) ([]T, []T, []T) {
	added := SetDifferenceByHash(to, from, hasher)
	removed := SetDifferenceByHash(from, to, hasher)
	unchanged := SetIntersectionByHash(from, to, hasher)
	return added, removed, unchanged
}

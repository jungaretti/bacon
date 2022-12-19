package collections

func SetDifferenceByHash[T any, U comparable](from, to []T, hasher func(T) U) []T {
	set := make(map[U]bool)

	for _, element := range to {
		set[hasher(element)] = true
	}

	var missing []T
	for _, element := range from {
		if _, ok := set[hasher(element)]; !ok {
			missing = append(missing, element)
		}
	}
	return missing
}

func DiffElementsByHash[T any, U comparable](from, to []T, hasher func(T) U) ([]T, []T) {
	added := SetDifferenceByHash(to, from, hasher)
	removed := SetDifferenceByHash(from, to, hasher)
	return added, removed
}

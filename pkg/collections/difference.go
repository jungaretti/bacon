package collections

func SetDifferenceByHash[T any](from, to []T, hasher func(T) string) []T {
	set := make(map[string]bool)

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

func DiffElementsByHash[T any](from, to []T, hasher func(T) string) ([]T, []T) {
	added := SetDifferenceByHash(to, from, hasher)
	removed := SetDifferenceByHash(from, to, hasher)
	return added, removed
}

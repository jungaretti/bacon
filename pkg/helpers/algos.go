package helpers

// Difference of two sets (A - B), which includes elements that are found in a but not in b
func Difference[T comparable](a, b []T) (diff []T) {
	m := make(map[T]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

// Difference of two sets (A - B), which includes elements that are found in a but not in b, with a custom comparator
func DifferenceByHasher[T any](a, b []T, hasher func(T) string) (diff []T) {
	m := make(map[string]bool)

	for _, item := range b {
		target := hasher(item)
		m[target] = true
	}

	for _, item := range a {
		target := hasher(item)
		if _, ok := m[target]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

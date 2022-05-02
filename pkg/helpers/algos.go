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

package slice_gotcha

// FilterVersions takes a slice of version strings and returns a new slice
// containing only versions that match the predicate. The original slice
// must NOT be modified.

func FilterVersions(versions []string, keep func(string) bool) []string {
	result := []string{}
	for _, v := range versions {
		if keep(v) {
			result = append(result, v)
		}
	}
	return result
}

// UniqueVersions returns a deduplicated copy of the input slice,
// preserving order of first occurrence. The input must NOT be modified.

func UniqueVersions(versions []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, v := range versions {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

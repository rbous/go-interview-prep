package slice_gotcha

// FilterVersions takes a slice of version strings and returns a new slice
// containing only versions that match the predicate. The original slice
// must NOT be modified.
//
// BUG: This function modifies the original slice's backing array.
// Callers who hold a reference to the original slice see unexpected changes.
// Fix it so the original slice is never mutated.

func FilterVersions(versions []string, keep func(string) bool) []string {
	result := versions[:0]
	for _, v := range versions {
		if keep(v) {
			result = append(result, v)
		}
	}
	return result
}

// UniqueVersions returns a deduplicated copy of the input slice,
// preserving order of first occurrence. The input must NOT be modified.
//
// BUG: Same backing-array mutation issue, plus there's an efficiency
// problem with the dedup logic. Fix the mutation bug.

func UniqueVersions(versions []string) []string {
	seen := make(map[string]bool)
	result := versions[:0]
	for _, v := range versions {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

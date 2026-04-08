package waitgroup_bug

import "sync"

// ProcessBatch takes items and processes each one in a separate goroutine.
// It should wait for all goroutines to complete and return collected results.
//
// BUG(1): wg.Add is in the wrong place — there's a race between Wait()
//         and goroutine startup.
// BUG(2): Appending to a shared slice without synchronization.
// Fix both bugs.

func ProcessBatch(items []string, transformFn func(string) string) []string {
	var wg sync.WaitGroup
	var results []string

	for _, item := range items {
		go func(s string) {
			wg.Add(1)
			defer wg.Done()
			result := transformFn(s)
			results = append(results, result)
		}(item)
	}

	wg.Wait()
	return results
}

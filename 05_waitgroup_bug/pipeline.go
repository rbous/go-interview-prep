package waitgroup_bug

import "sync"

// ProcessBatch takes items and processes each one in a separate goroutine.
// It should wait for all goroutines to complete and return collected results.

func ProcessBatch(items []string, transformFn func(string) string) []string {
	var wg sync.WaitGroup
	var results []string
	var mu sync.Mutex

	for _, item := range items {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			result := transformFn(s)
			mu.Lock()
			defer mu.Unlock()
			results = append(results, result)
		}(item)
	}

	wg.Wait()
	return results
}

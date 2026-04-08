package context_cancel

import (
	"context"
	"sync"
	"time"
)

// DownloadPackages simulates downloading packages concurrently.
// It should respect context cancellation: if ctx is cancelled,
// all in-flight downloads should stop promptly and the function
// should return whatever results were collected so far.

type DownloadResult struct {
	Name    string
	Size    int
	Success bool
}

func DownloadPackages(ctx context.Context, packages []string) []DownloadResult {
	var results []DownloadResult
	var wg sync.WaitGroup

	for _, pkg := range packages {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			size := simulateDownload(name)
			results = append(results, DownloadResult{
				Name:    name,
				Size:    size,
				Success: true,
			})
		}(pkg)
	}

	wg.Wait()
	return results
}

// simulateDownload pretends to download a package.
// It takes 500ms per package. This should be interruptible.
func simulateDownload(name string) int {
	time.Sleep(500 * time.Millisecond)
	return len(name) * 100
}

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

// add context cancellation-> select case <-ctx.Done()
// remove race condition   -> mutex

type DownloadResult struct {
	Name    string
	Size    int
	Success bool
}

func DownloadPackages(ctx context.Context, packages []string) []DownloadResult {
	var results []DownloadResult
	var wg sync.WaitGroup
	var mu sync.Mutex
	done := make(chan struct{})

	for _, pkg := range packages {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			size := simulateDownload(ctx, name)
			if size == -1 {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			results = append(results, DownloadResult{
				Name:    name,
				Size:    size,
				Success: true,
			})
		}(pkg)
	}

	go func(){
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}
	return results
}

// simulateDownload pretends to download a package.
// It takes 500ms per package. This should be interruptible.
func simulateDownload(ctx context.Context, name string) int {
	select {
	case <-time.After(500 * time.Millisecond):
		return len(name) * 100
	case <-ctx.Done():
		return -1
	} 
}

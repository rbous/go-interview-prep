package goroutine_leak

import (
	"errors"
	"runtime"
	"time"
)

// FetchAll takes a list of URLs and a fetch function, and returns results
// from all successful fetches. It should launch one goroutine per URL.

func FetchAll(urls []string, fetchFn func(string) (string, error)) []string {
	ch := make(chan string, len(urls))
	errCh := make(chan struct{}, len(urls))

	for _, url := range urls {
		go func(u string) {
			result, err := fetchFn(u)
			if err == nil {
				ch <- result
			} else {
				errCh <- struct{}{}
			}
		}(url)
	}

	var results []string
	timeout := time.After(2 * time.Second)
	for range urls {
		select {
		case r := <-ch:
			results = append(results, r)
		case <-errCh:
			continue
		case <-timeout:
			return results
		}
	}

	return results
}

// GoroutineCount returns the current number of goroutines.
func GoroutineCount() int {
	return runtime.NumGoroutine()
}

// SimulateFetch is a test helper that simulates a fetch.
// Odd-indexed URLs "fail", even-indexed ones succeed.
func SimulateFetch(url string) (string, error) {
	time.Sleep(10 * time.Millisecond)
	if url == "fail" {
		return "", errors.New("fetch failed")
	}
	return "data:" + url, nil
}

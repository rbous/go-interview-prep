package goroutine_leak

import (
	"errors"
	"runtime"
	"time"
)

// FetchAll takes a list of URLs and a fetch function, and returns results
// from all successful fetches. It should launch one goroutine per URL.
//
// BUG: Goroutines leak when not all results are consumed. After FetchAll
// returns, no goroutines from this call should still be blocked.
// Fix the goroutine leak.

func FetchAll(urls []string, fetchFn func(string) (string, error)) []string {
	ch := make(chan string)

	for _, url := range urls {
		go func(u string) {
			result, err := fetchFn(u)
			if err == nil {
				ch <- result
			}
		}(url)
	}

	var results []string
	timeout := time.After(2 * time.Second)
	for range urls {
		select {
		case r := <-ch:
			results = append(results, r)
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

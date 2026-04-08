package goroutine_leak

import (
	"testing"
	"time"
)

func TestFetchAllNoLeak(t *testing.T) {
	before := GoroutineCount()

	urls := []string{"a", "fail", "b", "fail", "c"}
	results := FetchAll(urls, SimulateFetch)

	// Give goroutines time to exit
	time.Sleep(200 * time.Millisecond)
	after := GoroutineCount()

	if len(results) != 3 {
		t.Errorf("got %d results, want 3", len(results))
	}

	// Allow +1 for test runtime goroutines
	if after > before+1 {
		t.Errorf("goroutine leak: before=%d, after=%d", before, after)
	}
}

func TestFetchAllEmpty(t *testing.T) {
	results := FetchAll(nil, SimulateFetch)
	if len(results) != 0 {
		t.Errorf("got %d results, want 0", len(results))
	}
}

func TestFetchAllAllSucceed(t *testing.T) {
	before := GoroutineCount()

	urls := []string{"x", "y", "z"}
	results := FetchAll(urls, SimulateFetch)

	time.Sleep(200 * time.Millisecond)
	after := GoroutineCount()

	if len(results) != 3 {
		t.Errorf("got %d results, want 3", len(results))
	}
	if after > before+1 {
		t.Errorf("goroutine leak: before=%d, after=%d", before, after)
	}
}

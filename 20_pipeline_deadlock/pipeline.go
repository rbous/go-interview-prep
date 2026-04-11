package pipeline_deadlock

import "sync"

// ProcessConcurrent applies process() to each item concurrently
// and returns all results (order does not matter).
//
// Each item is handled in its own goroutine. The caller provides
// the transformation function.
func ProcessConcurrent(items []string, process func(string) string) []string {
    results := make(chan string)
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        go func(s string) {
            defer wg.Done()
            results <- process(s)
        }(item)
    }

    wg.Wait()
    close(results)

    var out []string
    for r := range results {
        out = append(out, r)
    }
    return out
}

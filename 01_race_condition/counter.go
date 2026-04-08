package race_condition

import "sync"

// Counter tracks a shared count that can be incremented by multiple goroutines.
// IncrementConcurrently should spawn `n` goroutines, each incrementing the counter
// `perGoroutine` times. The final count should equal n * perGoroutine.
//
// BUG: This function has a data race. Fix it so the count is correct
// and `go test -race` passes.

type Counter struct {
	count int
	mutex sync.Mutex
	wg    sync.WaitGroup
}

func (c *Counter) Increment() {
	c.mutex.Lock()
	c.count++
	c.mutex.Unlock()
	c.wg.Done()
}

func (c *Counter) Value() int {
	return c.count
}

func IncrementConcurrently(n, perGoroutine int) int {
	c := &Counter{}

	for i := 0; i < n; i++ {
		c.wg.Add(perGoroutine)
		go func() {
			for j := 0; j < perGoroutine; j++ {
				c.Increment()
			}
		}()
	}
	c.wg.Wait()
	return c.Value()
}

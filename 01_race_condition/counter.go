package race_condition

// Counter tracks a shared count that can be incremented by multiple goroutines.
// IncrementConcurrently should spawn `n` goroutines, each incrementing the counter
// `perGoroutine` times. The final count should equal n * perGoroutine.
//
// BUG: This function has a data race. Fix it so the count is correct
// and `go test -race` passes.

type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) Value() int {
	return c.count
}

func IncrementConcurrently(n, perGoroutine int) int {
	c := &Counter{}

	for i := 0; i < n; i++ {
		go func() {
			for j := 0; j < perGoroutine; j++ {
				c.Increment()
			}
		}()
	}

	return c.Value()
}

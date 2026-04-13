package main

import (
	"fmt"
	"time"

	boss "go-interview-prep/26_final_boss"
)

func main() {
	fmt.Println("=== Final Boss: Throttled Task Dispatcher Debug ===")
	fmt.Println()

	// Case 1: Ordering — results should come back in original order
	fmt.Println("--- Case 1: Ordering (10 tasks, expect IDs 0-9 in order) ---")
	limiter1 := boss.NewTokenBucketLimiter(100)
	d1 := boss.NewDispatcher(limiter1, 5)

	taskCount := 10
	for i := 0; i < taskCount; i++ {
		d1.Submit(boss.Task{ID: i, Payload: fmt.Sprintf("msg-%d", i)})
	}

	results := d1.Results(taskCount)
	for i, res := range results {
		match := "✅"
		if res.TaskID != i {
			match = "❌ ORDER MISMATCH"
		}
		fmt.Printf("  results[%d]: TaskID=%d Output=%q  %s\n", i, res.TaskID, res.Output, match)
	}
	d1.Stop()
	fmt.Println()

	// Case 2: Count — d.count should equal taskCount after processing
	fmt.Println("--- Case 2: Count tracking (20 tasks, 10 workers) ---")
	limiter2 := boss.NewTokenBucketLimiter(1000)
	d2 := boss.NewDispatcher(limiter2, 10)

	for i := 0; i < 20; i++ {
		d2.Submit(boss.Task{ID: i, Payload: "data"})
	}
	_ = d2.Results(20)
	d2.Stop()
	// count is unexported, so we can't read it from main — but the race detector will catch it.
	// Run with: go run -race ./26_final_boss/cmd/
	fmt.Println("  (run with -race flag to check for data races on d.count)")
	fmt.Println()

	// Case 3: Graceful shutdown — Stop() should not hang or panic
	fmt.Println("--- Case 3: Graceful shutdown (slow rate limiter) ---")
	limiter3 := boss.NewTokenBucketLimiter(10) // slow
	d3 := boss.NewDispatcher(limiter3, 2)

	go func() {
		for i := 0; i < 5; i++ {
			d3.Submit(boss.Task{ID: i, Payload: "slow"})
		}
	}()

	time.Sleep(50 * time.Millisecond)

	done := make(chan struct{})
	go func() {
		d3.Stop()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("  ✅ Stop() returned cleanly")
	case <-time.After(2 * time.Second):
		fmt.Println("  ⚠️  Stop() blocked for >2s — potential deadlock!")
	}
	fmt.Println()

	// Case 4: WaitGroup — does wg.Done() ever get called?
	fmt.Println("--- Case 4: Submit 3 tasks, check if wg.Wait() ever unblocks ---")
	limiter4 := boss.NewTokenBucketLimiter(100)
	d4 := boss.NewDispatcher(limiter4, 2)

	for i := 0; i < 3; i++ {
		d4.Submit(boss.Task{ID: i, Payload: fmt.Sprintf("t-%d", i)})
	}
	// Drain results
	for i := 0; i < 3; i++ {
		res := d4.Results(1)
		fmt.Printf("  Got result: TaskID=%d\n", res[0].TaskID)
	}

	stopDone := make(chan struct{})
	go func() {
		d4.Stop()
		close(stopDone)
	}()

	select {
	case <-stopDone:
		fmt.Println("  ✅ Stop() returned — wg.Done() was called correctly")
	case <-time.After(2 * time.Second):
		fmt.Println("  ⚠️  Stop() hung — wg.Done() is probably never called!")
	}
}

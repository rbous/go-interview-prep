package final_boss

import (
	"fmt"
	"testing"
	"time"
)

func TestDispatcherOrderingAndConcurrency(t *testing.T) {
	limiter := NewTokenBucketLimiter(100)
	d := NewDispatcher(limiter, 10)

	taskCount := 20
	for i := 0; i < taskCount; i++ {
		d.Submit(Task{ID: i, Payload: fmt.Sprintf("msg-%d", i)})
	}

	results := d.Results(taskCount)

	// Verify order
	for i, res := range results {
		if res.TaskID != i {
			t.Errorf("Result %d has TaskID %d, want %d (Order Mismatch!)", i, res.TaskID, i)
		}
	}

	d.Stop()
}

func TestDispatcherRaceAndCount(t *testing.T) {
	limiter := NewTokenBucketLimiter(1000)
	d := NewDispatcher(limiter, 50)

	taskCount := 100
	for i := 0; i < taskCount; i++ {
		d.Submit(Task{ID: i, Payload: "data"})
	}

	_ = d.Results(taskCount)
	d.Stop()

	// d.count should be taskCount. 
	// This usually triggers -race if not synchronized.
	if d.count != taskCount {
		t.Errorf("Processed count = %d, want %d", d.count, taskCount)
	}
}

func TestGracefulShutdown(t *testing.T) {
	limiter := NewTokenBucketLimiter(10) // Slow rate
	d := NewDispatcher(limiter, 2)

	// Submit tasks that will take time due to rate limit
	go func() {
		for i := 0; i < 5; i++ {
			d.Submit(Task{ID: i, Payload: "slow"})
		}
	}()

	// Give it a tiny bit of time to start
	time.Sleep(50 * time.Millisecond)

	// Shutdown should not block forever and should not panic
	done := make(chan struct{})
	go func() {
		d.Stop()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Stop() timed out - potential deadlock")
	}
}

func TestNoLeaks(t *testing.T) {
	// This is a conceptual test. In a real environment, 
	// we'd check runtime.NumGoroutine()
}

package final_boss

import (
	"context"
	"sync"
	"time"
)

// The "Final Boss" Exercise: Throttled Task Dispatcher.
//
// This system accepts tasks, rate-limits them according to a strategy,
// processes them concurrently in a worker pool, and returns results
// in the original order. It must also support graceful shutdown.

type Task struct {
	ID      int
	Payload string
}

type Result struct {
	TaskID int
	Output string
	Err    error
}

// Limiter defines a rate-limiting strategy.
type Limiter interface {
	Wait(ctx context.Context) error
	Stop()
}

// TokenBucketLimiter allows 'rate' events per second.
type TokenBucketLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
}

func NewTokenBucketLimiter(rate int) *TokenBucketLimiter {
	l := &TokenBucketLimiter{
		tokens: make(chan struct{}, rate),
		ticker: time.NewTicker(time.Second / time.Duration(rate)),
	}
	// Background loop to fill the bucket.
	go func() {
		for range l.ticker.C {
			select {
			case l.tokens <- struct{}{}:
			default:
				// Bucket is full
			}
		}
	}()
	return l
}

func (l *TokenBucketLimiter) Wait(ctx context.Context) error {
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *TokenBucketLimiter) Stop() {
	// BUG 1: Forgot to stop the ticker! This causes a permanent background leak.
}

// Dispatcher manages throttled execution of tasks.
type Dispatcher struct {
	limiter Limiter
	tasks   chan Task
	results chan Result
	wg      sync.WaitGroup
	stopped bool
	count   int // Tracks total processed tasks
}

func NewDispatcher(l Limiter, workerCount int) *Dispatcher {
	d := &Dispatcher{
		limiter: l,
		tasks:   make(chan Task, 100),
		results: make(chan Result, 100),
	}

	for i := 0; i < workerCount; i++ {
		go d.worker()
	}

	return d
}

func (d *Dispatcher) worker() {
	for t := range d.tasks {
		// Respect rate limits
		d.limiter.Wait(context.Background())

		// Process the task
		res := Result{
			TaskID: t.ID,
			Output: "processed: " + t.Payload,
		}

		d.results <- res
		d.count++ // BUG 2: Data race on d.count (multiple workers writing)
	}
}

// Submit adds a task to the queue.
func (d *Dispatcher) Submit(t Task) {
	if d.stopped {
		return
	}
	d.wg.Add(1)
	d.tasks <- t
}

// Results collects all results in their original order.
func (d *Dispatcher) Results(count int) []Result {
	final := make([]Result, count)
	for i := 0; i < count; i++ {
		res := <-d.results
		// BUG 3: Ordering bug. Results are collected as they arrive, 
		// but they aren't placed back in the original index.
		final[i] = res
	}
	return final
}

// Stop gracefully shuts down the dispatcher.
func (d *Dispatcher) Stop() {
	d.stopped = true
	
	// BUG 4: Deadlock risk. If we close the channel before workers are done,
	// or wait in the wrong order, the system hangs.
	d.wg.Wait()
	close(d.tasks)
	
	// BUG 5: Forgot to stop the nested limiter.
}

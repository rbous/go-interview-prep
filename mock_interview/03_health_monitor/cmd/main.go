package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	hc "go-interview-prep/mock_interview/03_health_monitor"
)

// --- Debug mocks ---

type debugBootController struct {
	mu            sync.Mutex
	commitCalls   int
	rollbackCalls int
}

func (b *debugBootController) CommitCurrentSlot() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.commitCalls++
	fmt.Printf("  [BOOT] CommitCurrentSlot called (total: %d)\n", b.commitCalls)
	return nil
}

func (b *debugBootController) Rollback() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.rollbackCalls++
	fmt.Printf("  [BOOT] Rollback called (total: %d)\n", b.rollbackCalls)
	return nil
}

type debugCheck struct {
	name  string
	err   error
	delay time.Duration
}

func (c *debugCheck) Name() string { return c.name }
func (c *debugCheck) Check(ctx context.Context) error {
	fmt.Printf("  [CHECK] %s starting (delay: %v)\n", c.name, c.delay)
	select {
	case <-time.After(c.delay):
		if c.err != nil {
			fmt.Printf("  [CHECK] %s FAILED: %v\n", c.name, c.err)
		} else {
			fmt.Printf("  [CHECK] %s passed\n", c.name)
		}
		return c.err
	case <-ctx.Done():
		fmt.Printf("  [CHECK] %s cancelled\n", c.name)
		return ctx.Err()
	}
}

func main() {
	fmt.Println("=== Health Monitor Debug ===")
	fmt.Println()

	// Case 1: All checks pass — should commit, 0 rollbacks
	fmt.Println("--- Case 1: All checks pass (expect: 1 commit, 0 rollbacks) ---")
	boot1 := &debugBootController{}
	m1 := hc.NewMonitor(
		[]hc.HealthChecker{
			&debugCheck{name: "can_bus", delay: 10 * time.Millisecond},
			&debugCheck{name: "sensors", delay: 10 * time.Millisecond},
		},
		boot1,
		200*time.Millisecond,
	)
	err := m1.Run(context.Background())
	fmt.Printf("  Run returned: %v\n", err)
	// Wait for leaked watchdog
	time.Sleep(300 * time.Millisecond)
	boot1.mu.Lock()
	fmt.Printf("  Final: commits=%d, rollbacks=%d\n", boot1.commitCalls, boot1.rollbackCalls)
	boot1.mu.Unlock()
	fmt.Println()

	// Case 2: Check fails — should rollback ONCE
	fmt.Println("--- Case 2: Sensor fails (expect: 0 commits, 1 rollback) ---")
	boot2 := &debugBootController{}
	m2 := hc.NewMonitor(
		[]hc.HealthChecker{
			&debugCheck{name: "can_bus", delay: 10 * time.Millisecond},
			&debugCheck{name: "sensors", delay: 10 * time.Millisecond, err: fmt.Errorf("sensor offline")},
		},
		boot2,
		200*time.Millisecond,
	)
	err = m2.Run(context.Background())
	fmt.Printf("  Run returned: %v\n", err)
	time.Sleep(300 * time.Millisecond)
	boot2.mu.Lock()
	fmt.Printf("  Final: commits=%d, rollbacks=%d\n", boot2.commitCalls, boot2.rollbackCalls)
	boot2.mu.Unlock()
	fmt.Println()

	// Case 3: Slow check — watchdog should fire and cancel
	fmt.Println("--- Case 3: Slow check (expect: Run returns quickly, 1 rollback) ---")
	boot3 := &debugBootController{}
	m3 := hc.NewMonitor(
		[]hc.HealthChecker{
			&debugCheck{name: "slow_check", delay: 5 * time.Second},
		},
		boot3,
		100*time.Millisecond,
	)
	done := make(chan error)
	go func() { done <- m3.Run(context.Background()) }()
	select {
	case err := <-done:
		fmt.Printf("  Run returned: %v\n", err)
	case <-time.After(1 * time.Second):
		fmt.Println("  ⚠️  Run() blocked for >1s — watchdog didn't cancel pending checks!")
	}
	boot3.mu.Lock()
	fmt.Printf("  Final: commits=%d, rollbacks=%d\n", boot3.commitCalls, boot3.rollbackCalls)
	boot3.mu.Unlock()
	fmt.Println()

	// Case 4: Mutual exclusion — watchdog AND check both try to rollback
	fmt.Println("--- Case 4: Race — watchdog fires at 50ms, check fails at 150ms (expect: exactly 1 rollback) ---")
	boot4 := &debugBootController{}
	m4 := hc.NewMonitor(
		[]hc.HealthChecker{
			&debugCheck{name: "slow_fail", delay: 150 * time.Millisecond, err: fmt.Errorf("sensor fail")},
		},
		boot4,
		50*time.Millisecond,
	)
	done2 := make(chan error)
	go func() { done2 <- m4.Run(context.Background()) }()
	select {
	case err := <-done2:
		fmt.Printf("  Run returned: %v\n", err)
	case <-time.After(2 * time.Second):
		fmt.Println("  ⚠️  Run() blocked!")
	}
	time.Sleep(200 * time.Millisecond)
	boot4.mu.Lock()
	fmt.Printf("  Final: commits=%d, rollbacks=%d\n", boot4.commitCalls, boot4.rollbackCalls)
	boot4.mu.Unlock()
}

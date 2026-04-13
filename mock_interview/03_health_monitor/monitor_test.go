package healthcheck

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// --- Test mocks (do NOT modify) ---

type mockBootController struct {
	mu            sync.Mutex
	commitCalls   int
	rollbackCalls int
}

func (m *mockBootController) CommitCurrentSlot() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.commitCalls++
	return nil
}

func (m *mockBootController) Rollback() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rollbackCalls++
	return nil
}

type mockCheck struct {
	name  string
	err   error
	delay time.Duration
}

func (c *mockCheck) Name() string { return c.name }

func (c *mockCheck) Check(ctx context.Context) error {
	select {
	case <-time.After(c.delay):
		return c.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// --- Tests ---

func TestAllChecksPassCommits(t *testing.T) {
	boot := &mockBootController{}
	checks := []HealthChecker{
		&mockCheck{name: "can_bus", delay: 10 * time.Millisecond},
		&mockCheck{name: "sensors", delay: 10 * time.Millisecond},
	}

	m := NewMonitor(checks, boot, 200*time.Millisecond)
	err := m.Run(context.Background())
	if err != nil {
		t.Fatalf("all checks passed — expected commit, got error: %v", err)
	}

	// Wait for any leaked watchdog goroutine to fire.
	time.Sleep(300 * time.Millisecond)

	boot.mu.Lock()
	defer boot.mu.Unlock()
	if boot.commitCalls != 1 {
		t.Errorf("expected 1 commit, got %d", boot.commitCalls)
	}
	if boot.rollbackCalls != 0 {
		t.Errorf("expected 0 rollbacks, got %d — watchdog goroutine leaked and fired!", boot.rollbackCalls)
	}
}

func TestCheckFailureRollsBack(t *testing.T) {
	boot := &mockBootController{}
	checks := []HealthChecker{
		&mockCheck{name: "can_bus", delay: 10 * time.Millisecond},
		&mockCheck{name: "sensors", delay: 10 * time.Millisecond, err: fmt.Errorf("sensor offline")},
	}

	m := NewMonitor(checks, boot, 200*time.Millisecond)
	err := m.Run(context.Background())
	if err == nil {
		t.Fatal("expected error from failed health check")
	}

	// Wait for watchdog to potentially fire (double rollback).
	time.Sleep(300 * time.Millisecond)

	boot.mu.Lock()
	defer boot.mu.Unlock()
	if boot.rollbackCalls != 1 {
		t.Errorf("expected exactly 1 rollback, got %d", boot.rollbackCalls)
	}
	if boot.commitCalls != 0 {
		t.Errorf("expected 0 commits, got %d", boot.commitCalls)
	}
}

func TestWatchdogTriggersRollback(t *testing.T) {
	boot := &mockBootController{}
	checks := []HealthChecker{
		// This check is very slow — the watchdog should fire first
		// and cancel it, causing Run() to return promptly.
		&mockCheck{name: "slow_check", delay: 5 * time.Second},
	}

	m := NewMonitor(checks, boot, 100*time.Millisecond)

	done := make(chan error)
	go func() { done <- m.Run(context.Background()) }()

	select {
	case err := <-done:
		if err == nil {
			t.Error("expected error after watchdog timeout")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Run() did not return within 1s — watchdog should cancel pending health checks")
	}

	boot.mu.Lock()
	defer boot.mu.Unlock()
	if boot.rollbackCalls != 1 {
		t.Errorf("expected 1 rollback from watchdog, got %d", boot.rollbackCalls)
	}
	if boot.commitCalls != 0 {
		t.Errorf("expected 0 commits, got %d", boot.commitCalls)
	}
}

func TestRollbackMutualExclusion(t *testing.T) {
	boot := &mockBootController{}
	checks := []HealthChecker{
		// Watchdog fires at 50ms, check fails at 150ms.
		// Both try to rollback — only one should succeed.
		&mockCheck{name: "slow_fail", delay: 150 * time.Millisecond, err: fmt.Errorf("sensor fail")},
	}

	m := NewMonitor(checks, boot, 50*time.Millisecond)

	done := make(chan error)
	go func() { done <- m.Run(context.Background()) }()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run() did not return")
	}

	// Wait for everything to settle.
	time.Sleep(200 * time.Millisecond)

	boot.mu.Lock()
	defer boot.mu.Unlock()
	if boot.rollbackCalls != 1 {
		t.Errorf("expected exactly 1 rollback, got %d — commit/rollback must be mutually exclusive!", boot.rollbackCalls)
	}
}

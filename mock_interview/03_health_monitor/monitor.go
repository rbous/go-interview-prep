package healthcheck

import (
	"context"
	"fmt"
	"time"
)

// ========================================================================
// DESIGN DISCUSSION — Practice answering out loud (5-10 minutes):
//
// "After applying an OTA update, the vehicle reboots into the new A/B
// partition. How does the system determine whether the update is healthy
// and should be committed, or whether it should roll back to the
// previous partition?"
//
// Key points to discuss:
//   - Watchdog timer: if the system hasn't confirmed health within N
//     seconds, assume failure and reboot into the old partition.
//   - Health checks: CAN bus responsive, sensors online, critical
//     services started, dm-verity passes.
//   - Commit vs rollback must be mutually exclusive — once committed,
//     no rollback; once rolled back, no commit.
//   - Bootloader integration: U-Boot reads a boot_slot flag from a
//     small metadata partition. The watchdog resets this flag on failure.
//   - What if the health check process itself crashes? The hardware
//     watchdog should trigger a reboot into the old partition.
//   - Persisting the decision: the commit/rollback flag must survive
//     a reboot (written to disk, not just in memory).
//
// ========================================================================
//
// Now fix the buggy implementation below so all tests pass.
// Run: go test -race ./mock_interview/03_health_monitor/

// HealthChecker represents a single health check (e.g., "CAN bus responsive").
type HealthChecker interface {
	Name() string
	Check(ctx context.Context) error
}

// BootController manages the A/B partition boot slot.
type BootController interface {
	CommitCurrentSlot() error // Mark the current partition as good.
	Rollback() error         // Switch back to the previous partition.
}

// Monitor runs health checks after an update and decides whether to
// commit the new partition or roll back to the previous one.
type Monitor struct {
	checks    []HealthChecker
	boot      BootController
	timeout   time.Duration
	committed bool
}

func NewMonitor(checks []HealthChecker, boot BootController, timeout time.Duration) *Monitor {
	return &Monitor{
		checks:  checks,
		boot:    boot,
		timeout: timeout,
	}
}

// Run starts the post-update health verification process.
//
// It should:
//  1. Start a watchdog timer. If the timeout expires, roll back.
//  2. Run all health checks concurrently.
//  3. If ALL checks pass before the timeout, commit the current slot.
//  4. If ANY check fails OR the timeout expires, roll back.
//  5. Commit and rollback must be mutually exclusive (exactly one executes).
//  6. Run must return promptly after a decision is made (no lingering goroutines).
func (m *Monitor) Run(ctx context.Context) error {
	// Start the watchdog timer.
	go func() {
		time.Sleep(m.timeout)
		if !m.committed {
			m.boot.Rollback()
		}
	}()

	// Run all health checks concurrently.
	errs := make(chan error, len(m.checks))
	for _, check := range m.checks {
		go func(c HealthChecker) {
			errs <- c.Check(ctx)
		}(check)
	}

	// Collect results.
	for range m.checks {
		if err := <-errs; err != nil {
			m.boot.Rollback()
			return fmt.Errorf("health check failed: %w", err)
		}
	}

	// All checks passed — commit.
	m.committed = true
	return m.boot.CommitCurrentSlot()
}

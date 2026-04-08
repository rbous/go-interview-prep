package graceful_shutdown

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBasicSubmitAndShutdown(t *testing.T) {
	s := NewUpdateServer(2)

	s.Submit(Job{PackageName: "curl", Version: "8.0"})
	s.Submit(Job{PackageName: "wget", Version: "1.21"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := s.Shutdown(ctx)

	if len(results) != 2 {
		t.Errorf("got %d results, want 2", len(results))
	}
	for _, r := range results {
		if !r.Success {
			t.Errorf("package %s failed", r.PackageName)
		}
	}
}

func TestShutdownWaitsForInFlight(t *testing.T) {
	s := NewUpdateServer(2)

	for i := 0; i < 10; i++ {
		s.Submit(Job{PackageName: "pkg", Version: "1.0"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := s.Shutdown(ctx)

	if len(results) != 10 {
		t.Errorf("got %d results, want 10 (in-flight jobs should complete)", len(results))
	}
}

func TestRejectAfterShutdown(t *testing.T) {
	s := NewUpdateServer(2)
	s.Submit(Job{PackageName: "curl", Version: "8.0"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx)

	accepted := s.Submit(Job{PackageName: "late", Version: "1.0"})
	if accepted {
		t.Error("job was accepted after shutdown, should have been rejected")
	}
}

func TestConcurrentSubmit(t *testing.T) {
	s := NewUpdateServer(4)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Submit(Job{PackageName: "pkg", Version: "1.0"})
		}(i)
	}
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	results := s.Shutdown(ctx)

	if len(results) != 50 {
		t.Errorf("got %d results, want 50", len(results))
	}
}

func TestShutdownRespectsTimeout(t *testing.T) {
	s := NewUpdateServer(1)

	// Submit many slow jobs
	for i := 0; i < 100; i++ {
		s.Submit(Job{PackageName: "slow", Version: "1.0"})
	}

	// Very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	s.Shutdown(ctx)
	elapsed := time.Since(start)

	// Should not take more than ~500ms even with 100 jobs
	if elapsed > 1*time.Second {
		t.Errorf("Shutdown took %v; should respect context timeout", elapsed)
	}
}

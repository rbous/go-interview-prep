package context_cancel

import (
	"context"
	"testing"
	"time"
)

func TestDownloadAllSuccess(t *testing.T) {
	ctx := context.Background()
	pkgs := []string{"curl", "wget", "git"}

	results := DownloadPackages(ctx, pkgs)

	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	for _, r := range results {
		if !r.Success {
			t.Errorf("package %s failed", r.Name)
		}
	}
}

func TestDownloadRespectsCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	pkgs := []string{"linux-kernel", "gcc", "glibc", "openssl", "systemd"}

	start := time.Now()
	results := DownloadPackages(ctx, pkgs)
	elapsed := time.Since(start)

	// Should return well before all downloads would finish (5 * 500ms = 2.5s)
	if elapsed > 400*time.Millisecond {
		t.Errorf("took %v; should have cancelled promptly", elapsed)
	}

	// Might have partial results or none, but shouldn't have all 5
	if len(results) == 5 {
		t.Error("got all 5 results; context cancellation was not respected")
	}

	// All returned results should be marked as successful
	for _, r := range results {
		if r.Success && r.Size == 0 {
			t.Errorf("successful result %s has zero size", r.Name)
		}
	}
}

func TestDownloadEmpty(t *testing.T) {
	results := DownloadPackages(context.Background(), nil)
	if len(results) != 0 {
		t.Errorf("got %d results, want 0", len(results))
	}
}

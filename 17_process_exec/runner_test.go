package process_exec

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunCommandSuccess(t *testing.T) {
	stdout, stderr, err := RunCommand("echo", []string{"hello"}, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(stdout) != "hello" {
		t.Errorf("stdout = %q, want %q", stdout, "hello")
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}
}

func TestRunCommandFailure(t *testing.T) {
	_, stderr, err := RunCommand("sh", []string{"-c", "echo 'oops' >&2; exit 1"}, 5*time.Second)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(stderr, "oops") {
		t.Errorf("stderr should contain 'oops', got %q", stderr)
	}
	// The error message should include stderr for debugging
	if !strings.Contains(err.Error(), "oops") {
		t.Errorf("error should include stderr, got: %v", err)
	}
}

func TestRunCommandTimeout(t *testing.T) {
	start := time.Now()

	sleepCmd := "sleep"
	if runtime.GOOS == "windows" {
		t.Skip("sleep command differs on windows")
	}

	_, _, err := RunCommand(sleepCmd, []string{"30"}, 500*time.Millisecond)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed > 3*time.Second {
		t.Errorf("command took %v; timeout should have killed it at 500ms", elapsed)
	}
}

func TestRunScriptSuccess(t *testing.T) {
	out, err := RunScript("echo 'update complete'", 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "update complete") {
		t.Errorf("output = %q, want 'update complete'", out)
	}
}

func TestRunScriptTimeout(t *testing.T) {
	start := time.Now()
	_, err := RunScript("sleep 30", 500*time.Millisecond)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed > 3*time.Second {
		t.Errorf("script took %v; timeout should have killed it", elapsed)
	}
}

func TestRunWithRetrySuccess(t *testing.T) {
	ctx := context.Background()
	out, err := RunWithRetry(ctx, "echo", []string{"ok"}, 3, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "ok") {
		t.Errorf("output = %q, want 'ok'", out)
	}
}

func TestRunWithRetryAllFail(t *testing.T) {
	ctx := context.Background()
	_, err := RunWithRetry(ctx, "sh", []string{"-c", "exit 1"}, 3, 5*time.Second)
	if err == nil {
		t.Fatal("expected error after all retries")
	}
	if !strings.Contains(err.Error(), "3 attempts") {
		t.Errorf("error should mention attempt count: %v", err)
	}
}

func TestRunWithRetryRespectsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	start := time.Now()
	_, err := RunWithRetry(ctx, "sh", []string{"-c", "exit 1"}, 100, 5*time.Second)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error")
	}
	// Should not have tried 100 times — should have stopped after cancel
	if elapsed > 2*time.Second {
		t.Errorf("took %v; should have stopped promptly on context cancel", elapsed)
	}
}

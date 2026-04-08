package race_condition

import "testing"

func TestIncrementConcurrently(t *testing.T) {
	n := 100
	perGoroutine := 1000
	got := IncrementConcurrently(n, perGoroutine)
	want := n * perGoroutine

	if got != want {
		t.Errorf("IncrementConcurrently(%d, %d) = %d, want %d", n, perGoroutine, got, want)
	}
}

func TestIncrementConcurrentlySmall(t *testing.T) {
	got := IncrementConcurrently(5, 1)
	if got != 5 {
		t.Errorf("IncrementConcurrently(5, 1) = %d, want 5", got)
	}
}

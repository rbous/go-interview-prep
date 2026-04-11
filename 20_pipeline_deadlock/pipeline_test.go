package pipeline_deadlock

import (
    "sort"
    "strings"
    "testing"
)

func TestProcessConcurrent(t *testing.T) {
    items := []string{"alpha", "beta", "gamma", "delta"}
    got := ProcessConcurrent(items, strings.ToUpper)
    sort.Strings(got)
    want := []string{"ALPHA", "BETA", "DELTA", "GAMMA"}

    if len(got) != len(want) {
        t.Fatalf("got %d results, want %d: %v", len(got), len(want), got)
    }
    for i := range want {
        if got[i] != want[i] {
            t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
        }
    }
}

func TestProcessConcurrentSingle(t *testing.T) {
    got := ProcessConcurrent([]string{"hello"}, strings.ToUpper)
    if len(got) != 1 || got[0] != "HELLO" {
        t.Errorf("got %v, want [HELLO]", got)
    }
}

func TestProcessConcurrentEmpty(t *testing.T) {
    got := ProcessConcurrent([]string{}, strings.ToUpper)
    if len(got) != 0 {
        t.Errorf("got %v, want empty slice", got)
    }
}

package waitgroup_bug

import (
	"sort"
	"strings"
	"testing"
)

func TestProcessBatch(t *testing.T) {
	items := []string{"hello", "world", "foo", "bar"}
	results := ProcessBatch(items, strings.ToUpper)

	sort.Strings(results)
	want := []string{"BAR", "FOO", "HELLO", "WORLD"}

	if len(results) != len(want) {
		t.Fatalf("got %d results, want %d", len(results), len(want))
	}
	for i := range want {
		if results[i] != want[i] {
			t.Errorf("results[%d] = %q, want %q", i, results[i], want[i])
		}
	}
}

func TestProcessBatchEmpty(t *testing.T) {
	results := ProcessBatch(nil, strings.ToUpper)
	if len(results) != 0 {
		t.Errorf("got %d results, want 0", len(results))
	}
}

func TestProcessBatchSingle(t *testing.T) {
	results := ProcessBatch([]string{"test"}, strings.ToUpper)
	if len(results) != 1 || results[0] != "TEST" {
		t.Errorf("got %v, want [TEST]", results)
	}
}

func TestProcessBatchLarge(t *testing.T) {
	items := make([]string, 1000)
	for i := range items {
		items[i] = "item"
	}
	results := ProcessBatch(items, strings.ToUpper)
	if len(results) != 1000 {
		t.Errorf("got %d results, want 1000", len(results))
	}
}

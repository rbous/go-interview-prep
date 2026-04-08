package channel_panic

import (
	"sort"
	"testing"
)

func double(x int) int { return x * 2 }

func TestDispatch(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5}
	results := Dispatch(jobs, 3, double)

	sort.Ints(results)
	want := []int{2, 4, 6, 8, 10}

	if len(results) != len(want) {
		t.Fatalf("got %d results, want %d", len(results), len(want))
	}
	for i := range want {
		if results[i] != want[i] {
			t.Errorf("results[%d] = %d, want %d", i, results[i], want[i])
		}
	}
}

func TestDispatchEmpty(t *testing.T) {
	results := Dispatch(nil, 3, double)
	if len(results) != 0 {
		t.Errorf("got %d results, want 0", len(results))
	}
}

func TestDispatchOrdered(t *testing.T) {
	jobs := []int{10, 20, 30, 40, 50}
	results := DispatchOrdered(jobs, 3, double)

	want := []int{20, 40, 60, 80, 100}
	if len(results) != len(want) {
		t.Fatalf("got %d results, want %d", len(results), len(want))
	}
	for i := range want {
		if results[i] != want[i] {
			t.Errorf("results[%d] = %d, want %d", i, results[i], want[i])
		}
	}
}

func TestDispatchManyWorkers(t *testing.T) {
	jobs := []int{1, 2}
	results := Dispatch(jobs, 10, double)

	sort.Ints(results)
	if len(results) != 2 || results[0] != 2 || results[1] != 4 {
		t.Errorf("got %v, want [2 4]", results)
	}
}

package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestDependencyOrder(t *testing.T) {
	var mu sync.Mutex
	var order []string

	record := func(name string) UpdateFunc {
		return func(ctx context.Context) error {
			mu.Lock()
			order = append(order, name)
			mu.Unlock()
			return nil
		}
	}

	components := []Component{
		{Name: "gateway", Deps: nil, Update: record("gateway")},
		{Name: "bodyA", Deps: []string{"gateway"}, Update: record("bodyA")},
		{Name: "nav", Deps: nil, Update: record("nav")},
	}

	o := NewOrchestrator(components)

	done := make(chan []Result)
	go func() { done <- o.Execute(context.Background()) }()

	select {
	case results := <-done:
		for _, r := range results {
			if r.Err != nil {
				t.Errorf("%s failed: %v", r.Name, r.Err)
			}
		}
		mu.Lock()
		gatewayIdx, bodyAIdx := -1, -1
		for i, name := range order {
			if name == "gateway" {
				gatewayIdx = i
			}
			if name == "bodyA" {
				bodyAIdx = i
			}
		}
		mu.Unlock()
		if gatewayIdx > bodyAIdx {
			t.Error("gateway must execute before bodyA")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Execute timed out — possible deadlock")
	}
}

func TestFanOutDependency(t *testing.T) {
	// Two components depend on the same parent.
	// Both must be able to observe the parent's completion.
	components := []Component{
		{Name: "gateway", Update: func(ctx context.Context) error { return nil }},
		{Name: "bodyA", Deps: []string{"gateway"}, Update: func(ctx context.Context) error { return nil }},
		{Name: "bodyB", Deps: []string{"gateway"}, Update: func(ctx context.Context) error { return nil }},
	}

	o := NewOrchestrator(components)

	done := make(chan []Result)
	go func() { done <- o.Execute(context.Background()) }()

	select {
	case results := <-done:
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}
		for _, r := range results {
			if r.Err != nil {
				t.Errorf("%s should succeed: %v", r.Name, r.Err)
			}
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Execute deadlocked — hint: what happens when two goroutines read from a single-value buffered channel?")
	}
}

func TestConcurrentResultsNoRace(t *testing.T) {
	// 20 independent components — stresses the results collection.
	// Run with -race to catch data races.
	components := make([]Component, 20)
	for i := range components {
		name := fmt.Sprintf("comp-%d", i)
		components[i] = Component{
			Name:   name,
			Update: func(ctx context.Context) error { return nil },
		}
	}

	o := NewOrchestrator(components)

	done := make(chan []Result)
	go func() { done <- o.Execute(context.Background()) }()

	select {
	case results := <-done:
		if len(results) != 20 {
			t.Errorf("expected 20 results, got %d", len(results))
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out")
	}
}

func TestDependencyFailurePropagates(t *testing.T) {
	components := []Component{
		{Name: "gateway", Update: func(ctx context.Context) error {
			return fmt.Errorf("flash write failed")
		}},
		{Name: "body", Deps: []string{"gateway"}, Update: func(ctx context.Context) error {
			t.Error("body update should NOT run when gateway fails")
			return nil
		}},
	}

	o := NewOrchestrator(components)

	done := make(chan []Result)
	go func() { done <- o.Execute(context.Background()) }()

	select {
	case results := <-done:
		resultMap := make(map[string]error)
		for _, r := range results {
			resultMap[r.Name] = r.Err
		}
		if resultMap["gateway"] == nil {
			t.Error("gateway should have failed")
		}
		if resultMap["body"] == nil {
			t.Error("body should fail because its dependency (gateway) failed")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out")
	}
}

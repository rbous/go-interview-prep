package main

import (
	"context"
	"fmt"
	"time"

	orch "go-interview-prep/mock_interview/02_orchestrator"
)

func main() {
	fmt.Println("=== Component Orchestrator Debug ===")
	fmt.Println()

	// Case 1: Simple dependency chain — should work
	fmt.Println("--- Case 1: Simple chain (gateway -> body) ---")
	run([]orch.Component{
		{Name: "gateway", Update: makeUpdate("gateway", nil)},
		{Name: "body", Deps: []string{"gateway"}, Update: makeUpdate("body", nil)},
		{Name: "nav", Update: makeUpdate("nav", nil)},
	})

	// Case 2: Fan-out — two components depend on same parent
	// This is where the deadlock hides!
	fmt.Println("--- Case 2: Fan-out (gateway -> bodyA AND gateway -> bodyB) ---")
	run([]orch.Component{
		{Name: "gateway", Update: makeUpdate("gateway", nil)},
		{Name: "bodyA", Deps: []string{"gateway"}, Update: makeUpdate("bodyA", nil)},
		{Name: "bodyB", Deps: []string{"gateway"}, Update: makeUpdate("bodyB", nil)},
	})

	// Case 3: Dependency failure propagation
	fmt.Println("--- Case 3: Gateway fails — body should NOT run ---")
	run([]orch.Component{
		{Name: "gateway", Update: makeUpdate("gateway", fmt.Errorf("flash write failed"))},
		{Name: "body", Deps: []string{"gateway"}, Update: makeUpdate("body", nil)},
	})
}

func makeUpdate(name string, err error) orch.UpdateFunc {
	return func(ctx context.Context) error {
		fmt.Printf("  [RUNNING] %s\n", name)
		return err
	}
}

func run(components []orch.Component) {
	o := orch.NewOrchestrator(components)

	done := make(chan []orch.Result)
	go func() { done <- o.Execute(context.Background()) }()

	select {
	case results := <-done:
		fmt.Printf("  Got %d results:\n", len(results))
		for _, r := range results {
			if r.Err != nil {
				fmt.Printf("    ❌ %s: %v\n", r.Name, r.Err)
			} else {
				fmt.Printf("    ✅ %s: OK\n", r.Name)
			}
		}
	case <-time.After(3 * time.Second):
		fmt.Println("  ⚠️  DEADLOCK — Execute() did not return within 3s!")
	}
	fmt.Println()
}

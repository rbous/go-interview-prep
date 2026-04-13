package orchestrator

import (
	"context"
	"fmt"
	"sync"
)

// ========================================================================
// DESIGN DISCUSSION — Practice answering out loud (5-10 minutes):
//
// "A vehicle needs to update several components: a gateway ECU, two body
// controllers, and navigation maps. The body controllers depend on the
// gateway being updated first (it relays CAN messages). The nav maps
// have no dependencies.
//
// How would you design the orchestration layer in Go to run independent
// updates concurrently while respecting dependency ordering? What happens
// if the gateway update fails?"
//
// Key points to discuss:
//   - DAG representation (adjacency list) for dependency ordering
//   - Goroutine per component, channels or WaitGroups for signaling
//   - Independent components (nav maps) run in parallel with everything
//   - context.Context for cancellation when a critical component fails
//   - State machine per component: Pending → Running → Success | Failed
//   - Persisting state to disk in case of power loss mid-update
//   - What if two components depend on the same parent? (fan-out)
//
// ========================================================================
//
// Now fix the buggy implementation below so all tests pass.
// Run: go test -race ./mock_interview/02_orchestrator/

// UpdateFunc represents the work to update a single component.
type UpdateFunc func(ctx context.Context) error

// Component describes an updatable vehicle component and its dependencies.
type Component struct {
	Name   string
	Deps   []string   // names of components that must complete first
	Update UpdateFunc
}

// Result captures the outcome of a single component update.
type Result struct {
	Name string
	Err  error
}

// Orchestrator manages concurrent execution of component updates.
type Orchestrator struct {
	components map[string]Component
}

func NewOrchestrator(components []Component) *Orchestrator {
	m := make(map[string]Component)
	for _, c := range components {
		m[c.Name] = c
	}
	return &Orchestrator{components: m}
}

// Execute runs all component updates concurrently, respecting dependency order.
// A component must not start until all its dependencies have succeeded.
// If a dependency fails, the dependent component should also report failure
// without running its update function.
// Returns a result for every component.
func (o *Orchestrator) Execute(ctx context.Context) []Result {
	// Each component gets a channel that signals when it completes.
	done := make(map[string]chan error)
	for name := range o.components {
		done[name] = make(chan error, 1)
	}

	var results []Result
	var wg sync.WaitGroup

	for _, comp := range o.components {
		wg.Add(1)
		go func(c Component) {
			defer wg.Done()

			// Wait for all dependencies to complete.
			for _, dep := range c.Deps {
				select {
				case depErr := <-done[dep]:
					if depErr != nil {
						err := fmt.Errorf("%s: dependency %s failed: %w", c.Name, dep, depErr)
						results = append(results, Result{c.Name, err})
						done[c.Name] <- err
						return
					}
				case <-ctx.Done():
					results = append(results, Result{c.Name, ctx.Err()})
					done[c.Name] <- ctx.Err()
					return
				}
			}

			// All dependencies satisfied — run the update.
			err := c.Update(ctx)
			results = append(results, Result{c.Name, err})
			done[c.Name] <- err
		}(comp)
	}

	wg.Wait()
	return results
}

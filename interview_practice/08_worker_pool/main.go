package main

import (
    "fmt"
    "sync"
)

// Worker Pool
//
// A pool of workers processes tasks concurrently. Each task doubles its value.
// Results are returned as a slice indexed by task ID.
//
// Expected output:
//   results: [0 2 4 6 8]

type Task struct {
    ID    int
    Value int
}

func runPool(tasks []Task, workers int) []int {
    results := make([]int, len(tasks))
    var wg sync.WaitGroup
    jobs := make(chan Task, len(tasks))

    for i := 0; i < workers; i++ {
        go func() {
            for t := range jobs {
                results[len(results)-1-t.ID] = t.Value * 2
                wg.Done()
            }
        }()
    }

    for _, t := range tasks {
        wg.Add(1)
        jobs <- t
    }

    close(jobs)
    wg.Wait()
    return results
}

func main() {
    tasks := []Task{
        {0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4},
    }
    fmt.Println("results:", runPool(tasks, 3))
}

package channel_panic

import "sync"

// Dispatch sends jobs to `numWorkers` workers and collects results.
// Each worker applies `processFn` to its job and sends the result back.

func Dispatch(jobs []int, numWorkers int, processFn func(int) int) []int {
	jobCh := make(chan int)
	resultCh := make(chan int)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for job := range jobCh {
				resultCh <- processFn(job)
			}
			close(resultCh)
		}()
	}

	// Send jobs
	go func() {
		for _, j := range jobs {
			jobCh <- j
		}
		close(jobCh)
	}()

	// Collect results
	var results []int
	for r := range resultCh {
		results = append(results, r)
	}

	return results
}

// DispatchOrdered is like Dispatch but preserves input order.

type indexedResult struct {
	index int
	value int
}

func DispatchOrdered(jobs []int, numWorkers int, processFn func(int) int) []int {
	jobCh := make(chan indexedResult)
	resultCh := make(chan indexedResult)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				resultCh <- indexedResult{index: job.index, value: processFn(job.value)}
			}
		}()
	}

	go func() {
		for i, j := range jobs {
			jobCh <- indexedResult{index: i, value: j}
		}
		close(jobCh)
	}()

	var results = make([]int, len(jobs))
	for r := range resultCh {
		results[r.index] = r.value
	}

	return results
}

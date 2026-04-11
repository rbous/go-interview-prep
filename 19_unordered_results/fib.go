package unordered_results

// fib returns the nth Fibonacci number (0-indexed).
// fib(0)=0, fib(1)=1, fib(2)=1, fib(5)=5, fib(10)=55.
func fib(n int) int {
    if n <= 1 {
        return n
    }
    return fib(n-1) + fib(n-2)
}

// FibPipeline computes fib(inputs[i]) for each element concurrently
// and returns results in the same order as the input slice.
//
// Example: FibPipeline([]int{5, 0, 3}) should return []int{5, 0, 2}.
func FibPipeline(inputs []int) []int {
    results := make([]int, len(inputs))
    ch := make(chan int, len(inputs))

    for _, n := range inputs {
        go func(val int) {
            ch <- fib(val)
        }(n)
    }

    for i := range inputs {
        results[i] = <-ch
    }
    return results
}

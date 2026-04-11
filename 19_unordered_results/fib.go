package unordered_results

// fib returns the nth Fibonacci number (0-indexed).
// fib(0)=0, fib(1)=1, fib(2)=1, fib(5)=5, fib(10)=55.
func fib(n int) int {
    if n <= 1 {
        return n
    }
    return fib(n-1) + fib(n-2)
}

type indexed struct {
    index int
    value int
}

// FibPipeline computes fib(inputs[i]) for each element concurrently
// and returns results in the same order as the input slice.
//
// Example: FibPipeline([]int{5, 0, 3}) should return []int{5, 0, 2}.
func FibPipeline(inputs []int) []int {
    results := make([]int, len(inputs))
    ch := make(chan indexed, len(inputs))

    for i, n := range inputs {
        go func(key, val int) {
            ch <- indexed{key, fib(val)}
        }(i, n)
    }

    for range inputs {
        v := <-ch
        results[v.index] = v.value
    }
    return results
}

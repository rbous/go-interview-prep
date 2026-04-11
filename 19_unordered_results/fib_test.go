package unordered_results

import (
    "reflect"
    "testing"
)

func TestFibPipelineOrdered(t *testing.T) {
    tests := []struct {
        inputs []int
        want   []int
    }{
        {[]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 1, 2, 3, 5}},
        {[]int{5, 0, 3}, []int{5, 0, 2}},
        // fib(10)=55 is slow; fib(0) and fib(1) are instant.
        // Results arrive out-of-order unless the fix is in place.
        {[]int{10, 0, 1}, []int{55, 0, 1}},
        {[]int{8, 1, 0, 5}, []int{21, 1, 0, 5}},
    }

    for _, tt := range tests {
        got := FibPipeline(tt.inputs)
        if !reflect.DeepEqual(got, tt.want) {
            t.Errorf("FibPipeline(%v) = %v, want %v", tt.inputs, got, tt.want)
        }
    }
}

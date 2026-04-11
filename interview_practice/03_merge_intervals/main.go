package main

import (
	"fmt"
	"sort"
)

// Merge Intervals
//
// Given a collection of intervals, merge all overlapping intervals.
// Two intervals overlap if they share any point (including edges).
// For example, [1,2] and [2,3] overlap at point 2 and should merge to [1,3].
//
// Expected output:
//   merge([[1,3],[2,6],[8,10],[15,18]]) = [[1,6],[8,10],[15,18]]
//   merge([[1,4],[4,5]]) = [[1,5]]
//   merge([[1,2],[2,3],[3,4],[4,5]]) = [[1,5]]

func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return nil
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	merged := [][]int{intervals[0]}

	for _, curr := range intervals[1:] {
		last := merged[len(merged)-1]

		if curr[0] < last[1] {
			// Overlapping — extend the end
			if curr[1] > last[1] {
				last[1] = curr[1]
			}
		} else {
			// No overlap — add as new interval
			merged = append(merged, curr)
		}
	}

	return merged
}

func main() {
	fmt.Printf("merge([[1,3],[2,6],[8,10],[15,18]]) = %v\n",
		merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))

	fmt.Printf("merge([[1,4],[4,5]]) = %v\n",
		merge([][]int{{1, 4}, {4, 5}}))

	fmt.Printf("merge([[1,2],[2,3],[3,4],[4,5]]) = %v\n",
		merge([][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}))
}

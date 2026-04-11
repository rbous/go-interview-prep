package main

import "fmt"

// Two Sum
//
// Given an array of integers and a target, return the indices of the
// two numbers that add up to the target. Each input has exactly one
// solution, and you may not use the same element twice.
//
// Expected output:
//   twoSum([2, 7, 11, 15], 9) = [0, 1]
//   twoSum([3, 2, 4], 6) = [1, 2]
//   twoSum([3, 3], 6) = [0, 1]

func twoSum(nums []int, target int) [2]int {
	seen := make(map[int]int) // value -> index

	for i, n := range nums {
		seen[n] = i
		complement := target - n
		if j, ok := seen[complement]; ok && j != i {
			return [2]int{j, i}
		}
	}

	return [2]int{-1, -1}
}

func main() {
	fmt.Printf("twoSum([2, 7, 11, 15], 9) = %v\n", twoSum([]int{2, 7, 11, 15}, 9))
	fmt.Printf("twoSum([3, 2, 4], 6) = %v\n", twoSum([]int{3, 2, 4}, 6))
	fmt.Printf("twoSum([3, 3], 6) = %v\n", twoSum([]int{3, 3}, 6))
}

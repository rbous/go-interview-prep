package main

import "fmt"

// Longest Substring Without Repeating Characters
//
// Given a string, find the length of the longest substring that contains
// no duplicate characters. Uses the "sliding window" technique.
//
// Expected output:
//   lengthOfLongest("abcabcbb") = 3   (substring: "abc")
//   lengthOfLongest("bbbbb") = 1      (substring: "b")
//   lengthOfLongest("pwwkew") = 3     (substring: "wke")
//   lengthOfLongest("abba") = 2       (substring: "ab" or "ba")
//   lengthOfLongest("") = 0

func lengthOfLongest(s string) int {
	lastSeen := make(map[byte]int) // char -> last index where it was seen
	maxLen := 0
	left := 0

	for right := 0; right < len(s); right++ {
		ch := s[right]

		if idx, ok := lastSeen[ch]; ok && idx >= left{
			left = idx + 1
		}

		lastSeen[ch] = right

		windowLen := right - left + 1
		if windowLen > maxLen {
			maxLen = windowLen
		}
	}

	return maxLen
}

func main() {
	cases := []string{"abcabcbb", "bbbbb", "pwwkew", "abba", ""}
	for _, s := range cases {
		fmt.Printf("lengthOfLongest(%q) = %d\n", s, lengthOfLongest(s))
	}
}

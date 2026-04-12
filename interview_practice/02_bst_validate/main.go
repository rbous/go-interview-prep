package main

import (
	"fmt"
	"math"
)

// Validate Binary Search Tree
//
// A valid BST requires that for every node:
//   - ALL nodes in its left subtree have values LESS than the node
//   - ALL nodes in its right subtree have values GREATER than the node
//
// Expected output:
//   Tree 1 (valid BST):   true
//   Tree 2 (invalid BST): false
//   Tree 3 (tricky):      false

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func isValidBST(node *TreeNode) bool {
	return isBST(node, math.MinInt, math.MaxInt)
}

func isBST(node *TreeNode, smallest, biggest int) bool {
	if node == nil {
		return true
	}
	if node.Val <= smallest || node.Val >= biggest {
		return false
	}
	return isBST(node.Left, smallest, node.Val) && isBST(node.Right, node.Val, biggest)
}

func main() {
	// Tree 1:    2
	//           / \
	//          1   3
	tree1 := &TreeNode{Val: 2,
		Left:  &TreeNode{Val: 1},
		Right: &TreeNode{Val: 3},
	}

	// Tree 2:    5
	//           / \
	//          1   4
	//             / \
	//            3   6
	tree2 := &TreeNode{Val: 5,
		Left: &TreeNode{Val: 1},
		Right: &TreeNode{Val: 4,
			Left:  &TreeNode{Val: 3},
			Right: &TreeNode{Val: 6},
		},
	}

	// Tree 3:    5
	//           / \
	//          1   6
	//             / \
	//            3   7
	// Note: 3 is in the RIGHT subtree of 5, but 3 < 5. Invalid!
	tree3 := &TreeNode{Val: 5,
		Left: &TreeNode{Val: 1},
		Right: &TreeNode{Val: 6,
			Left:  &TreeNode{Val: 3},
			Right: &TreeNode{Val: 7},
		},
	}

	_ = math.MinInt64 // hint: might be useful

	fmt.Printf("Tree 1 (valid BST):   %v\n", isValidBST(tree1))
	fmt.Printf("Tree 2 (invalid BST): %v\n", isValidBST(tree2))
	fmt.Printf("Tree 3 (tricky):      %v\n", isValidBST(tree3))
}

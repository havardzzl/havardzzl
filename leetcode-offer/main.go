package leetcodeoffer

type ListNode struct {
	Val  int
	Next *ListNode
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func InorderTraversal(root *TreeNode, res *[]int) {
	if root == nil {
		return
	}
	InorderTraversal(root.Left, res)
	*res = append(*res, root.Val)
	InorderTraversal(root.Right, res)
}

func PreorderTraversal(root *TreeNode, res *[]int) {
	if root == nil {
		return
	}
	*res = append(*res, root.Val)
	PreorderTraversal(root.Left, res)
	PreorderTraversal(root.Right, res)
}

func PostorderTraversal(root *TreeNode, res *[]int) {
	if root == nil {
		return
	}
	PostorderTraversal(root.Left, res)
	PostorderTraversal(root.Right, res)
	*res = append(*res, root.Val)
}

var TestTree = &TreeNode{
	Val: 1,
	Left: &TreeNode{
		Val: 2,
		Left: &TreeNode{
			Val: 4,
		},
		Right: &TreeNode{
			Val: 5,
		},
	},
	Right: &TreeNode{
		Val: 3,
		Left: &TreeNode{
			Val: 6,
		},
		Right: &TreeNode{
			Val: 7,
		},
	},
}

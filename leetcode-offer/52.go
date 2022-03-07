package leetcodeoffer

func increasingBST(root *TreeNode) *TreeNode {
	dm := &TreeNode{}
	cur := dm
	var inorder func(n *TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)
		cur.Left = nil
		cur.Right = n
		cur = n
		inorder(n.Right)
	}
	inorder(root)
	return dm.Right
}

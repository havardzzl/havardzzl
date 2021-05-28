package leetcode

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return head
	}
	v := head.Val
	prev, cur := head, head.Next
	for cur != nil {
		if cur.Val == v {
			prev.Next = cur.Next
			cur = cur.Next
		} else {
			v = cur.Val
			prev, cur = cur, cur.Next
		}
	}
	return head
}

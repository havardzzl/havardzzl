package leetcode

func notDup(n1, n2 *ListNode) bool {
	return n1 == nil || n2 == nil || n1.Val != n2.Val
}

func deleteDuplicates2(head *ListNode) *ListNode {
	if head == nil {
		return head
	}
	var nHead, curN, prev *ListNode
	cur, next := head, head.Next
	for {
		if notDup(prev, cur) && notDup(cur, next) {
			if nHead == nil {
				nHead = cur
				curN = cur
			} else {
				curN.Next = cur
				curN = cur
			}
		}
		if next == nil {
			break
		}
		prev, cur, next = cur, next, next.Next
	}
	if curN != nil {
		curN.Next = nil
	}
	return nHead
}

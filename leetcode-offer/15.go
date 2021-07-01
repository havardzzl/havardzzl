package leetcodeoffer

func notDup(n1, n2 *ListNode) bool {
	return n1 == nil || n2 == nil || n1.Val != n2.Val
}

func HammingWeight(num uint32) int {
	var cnt int
	for num > 0 {
		if num%2 == 1 {
			cnt++
		}
		num /= 2
	}
	return cnt
}

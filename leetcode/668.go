package leetcode

// m乘n的乘法表，找到第k大的数 1<=k<=m*n

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// count of nums less than or equal to mid
func countByVal(m, n, mid int) int {
	var cnt int
	for i := 1; i <= m; i++ {
		cnt += min(n, mid/i)
	}
	return cnt
}

func inMatrix(m, n, v int) bool {
	for i := 1; i <= m; i++ {
		if v%i == 0 && v/i <= n {
			return true
		}
	}
	return false
}

func findKthNumber(m int, n int, k int) int {
	left, right := 1, k
	var mid int
	for left < right {
		mid = left + (right-left)/2
		cnt := countByVal(m, n, mid)
		switch {
		case cnt < k:
			left = mid + 1
		case cnt > k:
			right = mid
		default:
			for !inMatrix(m, n, mid) {
				mid--
			}
			return mid
		}
	}
	return right
}

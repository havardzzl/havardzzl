package leetcode

func findLengthOfLCIS(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	maxL := 1
	l, r := 0, 1
	for r < len(nums) {
		for r < len(nums) && nums[r] > nums[r-1] {
			r++
		}
		if maxL < r-l {
			maxL = r - l
		}
		l, r = r, r+1
	}
	return maxL
}

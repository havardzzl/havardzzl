package leetcode

func countTime(i int, manager []int, informTime []int) int {
	temp := manager[i]
	if temp == -1 {
		return 0
	}
	sumTime := informTime[temp]
	for manager[temp] != -1 {
		temp = manager[temp]
		sumTime += informTime[temp]
	}
	return sumTime
}

func numOfMinutes(n int, headID int, manager []int, informTime []int) int {
	var maxTime int
	counted := map[int]bool{}
	for i, timei := range informTime {
		if timei != 0 {
			continue
		}
		if counted[manager[i]] {
			continue
		}
		t := countTime(i, manager, informTime)
		if t > maxTime {
			maxTime = t
		}
		counted[manager[i]] = true
	}
	return maxTime
}

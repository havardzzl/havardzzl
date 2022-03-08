package leetcode

func largest1BorderedSquare(grid [][]int) int {
	row, col := len(grid), len(grid[0])
	max := row
	if col < max {
		max = col
	}
	check := func(grid [][]int, r0, c0, w int) bool {
		for c := c0; c < c0+w; c++ {
			if grid[r0][c] == 0 {
				return false
			}
		}
		for c := c0; c < c0+w; c++ {
			if grid[r0+w-1][c] == 0 {
				return false
			}
		}
		for r := r0; r < r0+w; r++ {
			if grid[r][c0] == 0 {
				return false
			}
		}
		for r := r0; r < r0+w; r++ {
			if grid[r][c0+w-1] == 0 {
				return false
			}
		}
		return true
	}
	for k := max; k > 0; k-- {
		for i := 0; i <= row-k; i++ {
			for j := 0; j <= col-k; j++ {
				if grid[i][j] == 1 && check(grid, i, j, k) {
					return k * k
				}
			}
		}
	}
	return 0
}

package main

import "fmt"

func gMap[T1 any, T2 any](arr []T1, f func(T1) T2) []T2 {
	result := make([]T2, len(arr))
	for i, elem := range arr {
		result[i] = f(elem)
	}
	return result
}

func main() {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	squares := gMap(nums, func(elem int) int {
		return elem * elem
	})
	fmt.Println(squares)
}

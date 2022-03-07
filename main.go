package main

import (
	"fmt"

	leetcodeoffer "github.com/havardzzl/havardzzl/leetcode-offer"
)

func main() {
	res := &[]int{}
	leetcodeoffer.InorderTraversal(leetcodeoffer.TestTree, res)
	fmt.Println(*res)
}

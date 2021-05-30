package leetcode

import "fmt"

// 思路：先预处理，预处理结果用两个int到string的map存储，如果前左（右）边i个字符能组成回文，i映射到余下字符逆序组成的
// 字符串

// 左侧前i个字符预处理
func prepareLeft(word string, leftI int) (ispalindrome bool, rightPart string) {
	l, r := 0, leftI-1
	for l < r {
		if word[l] != word[r] {
			return false, ""
		}
		l++
		r--
	}
	rightBs := make([]byte, len(word)-leftI)
	for i := 0; i < len(rightBs)/2; i++ {
		rightBs[i], rightBs[len(rightBs)-1-i] = word[len(word)-i], word[leftI+i]
	}
	return true, string(rightBs)
}

// 右侧后i个字符预处理
func prepareRight(word string, rightI int) (ispalindrome bool, leftPart string) {
	l, r := len(word)-rightI, len(word)-1
	for l < r {
		if word[l] != word[r] {
			return false, ""
		}
		l++
		r--
	}
	leftBs := make([]byte, len(word)-rightI)
	for i := 0; i < len(leftBs)/2; i++ {
		leftBs[i], leftBs[len(leftBs)-1-i] = word[len(leftBs)-1-i], word[i]
	}
	return true, string(leftBs)
}

func palindromePairs(words []string) [][]int {
	wordsL := map[int]map[int]string{}
	wordsR := map[int]map[int]string{}
	for wordi, word := range words {
		wordsL[wordi] = map[int]string{}
		wordsR[wordi] = map[int]string{}
		for i := range word {
			ispalindromeL, rp := prepareLeft(word, i+1)
			ispalindromeR, lp := prepareRight(word, i+1)
			if ispalindromeL {
				wordsL[wordi][i+1] = rp
			}
			if ispalindromeR {
				wordsR[wordi][i+1] = lp
			}
		}
	}
	fmt.Println(wordsL)
	fmt.Println(wordsR)
	return nil
}

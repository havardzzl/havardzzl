package leetcode

var alphm = map[byte]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true,
	'A': true, 'E': true, 'I': true, 'O': true, 'U': true}

func isOlph(a byte) bool {
	return alphm[a]
}

func reverseVowels(s string) string {
	l, r := 0, len(s)-1
	ns := []byte(s)
	for l < r {
		if !isOlph(s[l]) {
			l++
			continue
		}
		if !isOlph(s[r]) {
			r--
			continue
		}
		ns[l], ns[r] = ns[r], ns[l]
		l++
		r--
	}
	return string(ns)
}

package main

import (
	"fmt"
	"time"
	// "github.com/havardzzl/havardzzl/leetcodeoffer"
)

type K struct {
	B int
	A int
}

func (k *K) Chdl() {
	fmt.Println("chdl")
}

func kdf(k *K) {
	defer updateMutationElapsedTimeMetrics()()
	time.Sleep(time.Second)
	k.Chdl()
}

func main() {
	kdf(nil)
}

func updateMutationElapsedTimeMetrics() func() {
	start := time.Now()
	return func() {
		end := time.Now()
		fmt.Println("Elapsed time:", end.Sub(start))
	}
}

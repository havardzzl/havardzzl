package main

import (
	"fmt"
	"sync"

	"github.com/havardzzl/havardzzl/db"
)

var (
	cc = make(chan int)
	wg sync.WaitGroup
)

func workCc(cnt int) {
	defer wg.Done()
	fmt.Println("I'am wait for cc: ", cnt)
	for {
		i, ok := <-cc
		if !ok {
			fmt.Println("cc closed quit")
			return
		}
		fmt.Println("get from cc:", i, " I'am ", cnt)
	}
}

func main() {
	db.Init()
	db.NewTask()
	db.Close()

	// con := 3
	// for i := 0; i < con; i++ {
	// 	go workCc(i)
	// }
	// wg.Add(con)
	// go func() {
	// 	time.Sleep(time.Second)
	// 	cc <- 5
	// 	time.Sleep(time.Second * 2)
	// 	cc <- 6
	// 	time.Sleep(time.Second * 3)
	// 	cc <- 7
	// 	close(cc)
	// }()
	// wg.Wait()
}

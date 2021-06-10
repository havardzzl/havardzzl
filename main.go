package main

import (
	"fmt"
	"os"
	"time"

	"github.com/havardzzl/havardzzl/leetcode"
)

func writeLogToFile(msg string) {
	m := time.Now().Minute()
	filename := fmt.Sprintf("/Users/admin/go/src/z/havardzzl/%d.txt", m)
	envFilename := os.Getenv("log_file_name")
	if envFilename != "" {
		filename = envFilename
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("open log file error: %v", err)
		return
	}
	defer f.Close()
	if _, err = f.WriteString(msg); err != nil {
		panic(err)
	}
}

func tryPa() bool {
	defer func() {
		recover()
		// return "from panic"
	}()
	var k []int
	k[33] = 3
	return true
}

func main() {
	fmt.Println(leetcode.FindKthNumber(45, 12, 471))
	// fmt.Println(leetcode.FindKthNumber(9895, 28405, 100787757))

	// writeLogToFile("this is first log\n")
	// writeLogToFile("this is second log\n")
	// writeLogToFile("this is thrid log\n")
}

package main

import (
	"fmt"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":9092", http.FileServer(http.Dir("../assets")))
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}

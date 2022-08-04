package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	bs, _ := ioutil.ReadAll(r.Body)
	fmt.Println("get body: " + string(bs))
	r.Body.Close()
	fmt.Fprintf(w, "Hello World")
}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	ms := 100
	msr := r.URL.Query().Get("ms")
	if msr != "" {
		ms, _ = strconv.Atoi(msr)
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
	fmt.Fprintf(w, "Hi")
}

func main() {
	http.HandleFunc("/", HelloHandler)
	http.HandleFunc("/time", TimeHandler)
	fmt.Println("begin listening at :8979")
	http.ListenAndServe(":8979", nil)
}

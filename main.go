package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	bs, _ := ioutil.ReadAll(r.Body)
	fmt.Println("get body: " + string(bs))
	r.Body.Close()
	fmt.Fprintf(w, "Hello World")
}

func main() {
	http.HandleFunc("/", HelloHandler)
	fmt.Println("begin listening at :8979")
	http.ListenAndServe(":8979", nil)
}

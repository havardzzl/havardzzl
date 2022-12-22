package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/havardzzl/havardzzl/probe"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	bs, _ := ioutil.ReadAll(r.Body)
	fmt.Println("get body: " + string(bs))
	r.Body.Close()
	fmt.Fprintf(w, "Hello World")
}

func KubeHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("kubectl", "get", "nodes")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, string(output))
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
	// http.HandleFunc("/", HelloHandler)
	// http.HandleFunc("/kube", KubeHandler)
	// http.HandleFunc("/time", TimeHandler)
	// fmt.Println("begin listening at :8979")
	// http.ListenAndServe(":8979", nil)

	fmt.Println("probe result:", probe.ProbeTls())
}

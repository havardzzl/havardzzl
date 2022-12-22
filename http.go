package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"
	_ "time/tzdata"
)

const (
	con = 40
	cnt = math.MaxInt
	qps = 100
)

var (
	wg      sync.WaitGroup
	httpCli = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 3 * time.Second,
	}
	totalReq = 0
	failReq  = 0
)

func request(url string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := httpCli.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp == nil {
		fmt.Println("resp is nil")
		return
	}
	totalReq++
	if totalReq%1000 == 0 {
		fmt.Println("total request: ", totalReq)
	}
	if resp.StatusCode != 200 {
		failReq++
		bs := []byte{}
		if resp != nil && resp.Body != nil {
			bs, _ = io.ReadAll(resp.Body)
		}
		fmt.Printf("get status code %d, msg: %s, total: %d, fail: %d\n", resp.StatusCode, bs, totalReq, failReq)
	} else if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

}

func httpWorker(url string) {
	for i := 0; i < cnt; i++ {
		request(url)
		// time.Sleep(time.Second / qps)
	}
	wg.Done()
}

func GetMd5Str(str string, size int) string {
	md5time := md5.Sum([]byte(str))
	if size == 0 {
		size = 1
	} else if size > 16 {
		size = 16
	}
	return hex.EncodeToString(md5time[0:size])
}

func mainH() {
	wg.Add(con)
	for i := 0; i < con; i++ {
		go httpWorker("http://py-hello-z.serverless.kanzhun.tech")
	}
	wg.Wait()
}

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	period        = 100 * time.Millisecond
	cpuWorkTimes  = 2
	scheduleTimes = 100
)

// 大概两毫秒
func cal() {
	w := md5.New()
	for i := 0; i < 10000; i++ {
		ts := time.Now().UnixNano()
		s := fmt.Sprintf("%d", ts)
		io.WriteString(w, s)
	}
}

func cpuWork(period time.Duration, num int) {
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			var start, total int64
			for {
				if total > int64(period) {
					break
				}
				start = time.Now().UnixNano()
				cal()
				total += time.Now().UnixNano() - start
			}
		}()
	}
	wg.Wait()
}

// 给定cpu request/limit 1c
// 没有cpu burst时，理论上耗时：100 * (100ms + 2 * 100ms) = 30s
// 有cpu burst时，理论上耗时：100 * (100ms + 100ms) = 20s
func schedule() {
	for i := 0; i < scheduleTimes; i++ {
		time.Sleep(period)
		cpuWork(period, cpuWorkTimes)
	}
}

func test1() string {
	res := ""
	start := time.Now()
	cal()
	cost := time.Since(start)
	res += fmt.Sprintf("each cal cost: %s", cost)
	fmt.Println("each cal cost: ", cost)

	start = time.Now()
	cpuWork(period, cpuWorkTimes)
	cost = time.Since(start)
	res += fmt.Sprintf(", each cpu work cost: %s", cost)
	fmt.Println("each cpuWork cost: ", cost)

	start = time.Now()
	schedule()
	cost = time.Since(start)
	res += fmt.Sprintf(", each schedule cost: %s", cost)
	fmt.Println("total cost: ", cost)
	return res
}

func Routes(r *gin.Engine) {
	r.GET("/api/go_git", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Welcome to Gin Framework",
			"headers:": c.Request.Header,
		})
	})
	r.GET("/set", func(c *gin.Context) {
		p := struct {
			CpuWorkTimes  int `form:"cpuWorkTimes"`
			ScheduleTimes int `form:"scheduleTimes"`
		}{}
		if c.Bind(&p) == nil {
			if p.CpuWorkTimes > 0 {
				cpuWorkTimes = p.CpuWorkTimes
			}
			if p.ScheduleTimes > 0 {
				scheduleTimes = p.ScheduleTimes
			}
			c.String(http.StatusOK, "set cpuWorkTimes to %d, scheduleTimes to %d", cpuWorkTimes, scheduleTimes)
		}
	})
	r.GET("/run", func(c *gin.Context) {
		res := test1()
		c.String(http.StatusOK, res)
	})
}

func main() {
	runtime.GOMAXPROCS(20)
	var handler http.Handler
	router := gin.Default()
	// router.Use(UserMiddleware)
	Routes(router)
	handler = router
	srv := &http.Server{
		Addr:         ":8386",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	srv.Serve(listener)
}

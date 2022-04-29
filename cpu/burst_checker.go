package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	maxCalTime    time.Duration
	minCalTime    time.Duration = time.Hour
	avgCalTime    time.Duration
	period        = 200 * time.Millisecond
	cpuWorkTimes  = 6
	scheduleTimes = 10
)

func cal() {
	var k int64 = 3759304
	var i int64
	for i = 1; i < k; i++ {
		if k%i == 0 {
			i = i + 1
		}
	}
}

func cpuWork(calTime time.Duration, workers int) {
	times := int(calTime / avgCalTime)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			var time int
			for {
				cal()
				time++
				if time >= times {
					break
				}
			}
		}()
	}
	wg.Wait()
}

func schedule() {
	for i := 0; i < scheduleTimes; i++ {
		time.Sleep(2 * period)
		cpuWork(period, cpuWorkTimes)
	}
}

func test1() string {
	res := ""
	start := time.Now()
	cpuWork(period, cpuWorkTimes)
	cost := time.Since(start)
	fmt.Println("each cpuWork cost: ", cost)

	start = time.Now()
	schedule()
	cost = time.Since(start)
	res += fmt.Sprintf(", schedule cost: %s", cost)
	fmt.Println("total cost: ", cost)
	return res
}

func Routes(r *gin.Engine) {
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
	var start time.Time
	var total time.Duration
	calTime := 100
	for i := 0; i < calTime; i++ {
		start = time.Now()
		cal()
		cost := time.Since(start)
		total += cost
		if cost > maxCalTime {
			maxCalTime = cost
		}
		if cost < minCalTime {
			minCalTime = cost
		}
	}
	avgCalTime = total / time.Duration(calTime)
	fmt.Printf("cal cost: max: %s, min: %s, avg: %s\n", maxCalTime, minCalTime, avgCalTime)

	runtime.GOMAXPROCS(16)
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

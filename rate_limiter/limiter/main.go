package main

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// 每秒10个令牌  桶容量5
	limiter := rate.NewLimiter(10, 5)
	cnt := 0
	last := time.Now()
	for {
		ok := limiter.Allow()
		if ok {
			cur := time.Now()
			cnt++
			fmt.Println(cnt, cur.Sub(last))
			last = cur
		} else {
			time.Sleep(time.Microsecond * 20)
			// fmt.Println("reach limit, slow down")
		}
	}
}

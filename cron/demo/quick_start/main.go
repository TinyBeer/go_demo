package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	start := time.Now()
	// 使用带秒的解释器时等价于 * * * * * *
	c.AddFunc("@every 1s", func() {
		fmt.Println(time.Since(start))
	})

	c.Start()

	time.Sleep(time.Second * 10)
}

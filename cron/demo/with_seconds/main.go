package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())
	cnt := 0
	c.AddFunc("10-30/5 * * * * *", func() {
		cnt++
		fmt.Println(cnt, time.Now())
	})
	c.Start()

	time.Sleep(time.Minute * 3)

}

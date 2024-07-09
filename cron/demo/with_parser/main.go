package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	// 仅关注每天内的任务 并且支持 秒字段可选
	c := cron.New(cron.WithParser(cron.NewParser(cron.Second | cron.SecondOptional | cron.Minute | cron.Hour)))
	// 每分钟的第3秒执行
	c.AddFunc("3 * *", func() {
		fmt.Println(time.Now())
	})

	// 每分钟第0秒执行
	c.AddFunc("* *", func() {
		fmt.Println(time.Now())
	})

	c.Start()

	time.Sleep(time.Minute * 10)
}

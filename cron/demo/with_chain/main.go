package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	// c := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	c := cron.New(cron.WithChain(cron.DelayIfStillRunning(cron.DefaultLogger)))
	c.AddFunc("@every 1s", func() {
		time.Sleep(time.Second * 2)
		fmt.Println(time.Now())
	})

	c.AddJob(
		"@every 1s",
		cron.NewChain(
			cron.DelayIfStillRunning(cron.DefaultLogger),
		).Then(
			cron.FuncJob(func() {
				fmt.Println("hahha")
				time.Sleep(time.Second * 2)
			}),
		),
	)
	c.Start()

	time.Sleep(time.Second * 10)
}

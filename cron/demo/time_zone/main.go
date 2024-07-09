package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	c := cron.New(cron.WithLocation(loc))
	c.AddFunc("* * * * *", func() {
		fmt.Println(c.Location())
	})

	c.AddFunc("CRON_TZ=Asia/Tokyo 13 11 * * *", func() {
		loc, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			panic(err)
		}
		fmt.Println("tokyo time", time.Now().In(loc))
	})
	c.Start()

	time.Sleep(time.Second * 50)
}

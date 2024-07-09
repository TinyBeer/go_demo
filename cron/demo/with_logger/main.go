package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron:", log.LstdFlags))))
	c.AddFunc("@every 2s", func() {
		fmt.Println(time.Now())
	})
	c.Start()
	time.Sleep(time.Second * 10)
}

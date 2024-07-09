package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	myjob := new(MyJob)
	c.AddJob("@every 1s", myjob)
	c.AddJob("@every 2s", myjob)
	c.Start()

	time.Sleep(time.Second * 10)
}

type MyJob struct {
	mu sync.Mutex
	id int
}

func (m *MyJob) Run() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.id++
	fmt.Println("work with lock", m.id, time.Now())
	time.Sleep(time.Second * 3)
	fmt.Println("finish", m.id, time.Now())
}

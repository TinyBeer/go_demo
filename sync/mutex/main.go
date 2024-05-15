package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	l  sync.Mutex
	rl sync.RWMutex
)

func main() {
	LockTest()
	RLockTest()

	time.Sleep(time.Second * 20)
}

func LockTest() {
	go func() {
		l.Lock()
		defer l.Unlock()
		fmt.Println("locked 1", time.Now())
		time.Sleep(time.Second * 5)
	}()
	go func() {
		l.Lock()
		defer l.Unlock()
		fmt.Println("locked 2", time.Now())
		time.Sleep(time.Second * 5)
	}()
}

func RLockTest() {
	go func() {
		time.Sleep(0)
		rl.Lock()
		defer rl.Unlock()
		fmt.Println("rlocked 3", time.Now())
		time.Sleep(time.Second * 5)
	}()
	go func() {
		rl.RLock()
		defer rl.RUnlock()
		fmt.Println("rlocked 1", time.Now())
		time.Sleep(time.Second * 5)
	}()
	go func() {
		rl.RLock()
		defer rl.RUnlock()
		fmt.Println("rlocked 2", time.Now())
		time.Sleep(time.Second * 5)
	}()
}

package main

import (
	"log"
	"sync"
	"time"
)

var done = false

func main() {
	var l sync.Mutex
	c := sync.NewCond(&l)
	go read("reader one", c)
	go read("reader two", c)
	go read("reader three", c)
	go write("writer one", c)

	time.Sleep(time.Second * 3)
}

func read(name string, c *sync.Cond) {
	c.L.Lock()
	defer c.L.Unlock()
	for !done {
		c.Wait()
	}
	log.Println(name, "start reading...")
}

func write(name string, c *sync.Cond) {
	log.Println(name, "start writing...")
	time.Sleep(time.Second)
	c.L.Lock()
	done = true
	c.L.Unlock()
	log.Println("waik up all")
	c.Broadcast()

}

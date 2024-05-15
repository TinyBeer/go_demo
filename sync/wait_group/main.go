package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)
			fmt.Printf("mission %d complete!\n", id)
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("wati done")
}

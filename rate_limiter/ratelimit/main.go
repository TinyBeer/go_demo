package main

import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"
)

func main() {
	rl := ratelimit.New(10) // per second

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}

	// Output:
	// 0 0
	// 1 100ms
	// 2 100ms
	// 3 100ms
	// 4 100ms
	// 5 100ms
	// 6 100ms
	// 7 100ms
	// 8 100ms
	// 9 100ms
}

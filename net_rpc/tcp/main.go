package main

import "time"

func main() {
	go server()

	time.Sleep(time.Second)

	client()
}

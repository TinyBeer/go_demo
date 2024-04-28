package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	l := &lumberjack.Logger{
		Filename:   "./foo.log",
		MaxSize:    1,
		MaxAge:     1,
		MaxBackups: 2,
		LocalTime:  false,
		Compress:   true,
	}
	log.SetOutput(l)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	go func() {
		<-c
		fmt.Println("rotate")
		l.Rotate()
		os.Exit(0)
	}()

	cnt := 0
	for {
		time.Sleep(time.Millisecond * 10)
		log.Println(strings.Repeat("I", cnt))
		if cnt < 10000 {
			cnt++
		}
	}
}

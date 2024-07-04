package main

import (
	"fmt"
	"learn_xorm/common"
	"time"
)

func main() {
	e := common.Engine()
	for {
		err := e.Ping()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("healthy")
		}
		time.Sleep(time.Second)
	}
}

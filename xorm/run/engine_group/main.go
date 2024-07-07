package main

import (
	"fmt"
	"learn_xorm/common"
)

func main() {
	eg := common.NewEngineGroup()
	err := eg.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("ok")
}

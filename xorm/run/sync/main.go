package main

import (
	"fmt"
	"learn_xorm/common"
	"learn_xorm/model"
)

func main() {
	e := common.Engine()
	// _, err := e.SyncWithOptions(xorm.SyncOptions{WarnIfDatabaseColumnMissed: true}, &model.User{})
	err := e.Sync(&model.User{})
	if err != nil {
		panic(err)
	}
	fmt.Println("ok")
}

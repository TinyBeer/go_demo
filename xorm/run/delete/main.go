package main

import (
	"fmt"
	"learn_xorm/common"
	"learn_xorm/model"
)

func main() {
	e := common.Engine()
	affect, err := e.ID(6).Delete(&model.User{})
	if err != nil {
		panic(err)
	}
	fmt.Println("affected", affect)
}

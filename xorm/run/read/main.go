package main

import (
	"fmt"
	"learn_xorm/common"
	"learn_xorm/model"
)

func main() {
	e := common.Engine()
	u := &model.User{}
	ok, err := e.ID(2).Get(u)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("not exists")
		return
	}
	fmt.Println(u)

	var us []*model.User
	err = e.Find(&us)
	if err != nil {
		panic(err)
	}
	for _, u := range us {
		fmt.Println(u)
	}
}

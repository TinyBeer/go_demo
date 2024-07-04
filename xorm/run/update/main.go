package main

import (
	"fmt"
	"learn_xorm/common"
	"learn_xorm/model"
)

func main() {
	e := common.Engine()

	affect, err := e.ID(1).Update(&model.User{Name: "jack"})
	if err != nil {
		panic(err)
	}
	fmt.Println("affect row", affect)

	affect, err = e.Table(&model.User{}).ID(2).Update(map[string]interface{}{"age": 16})
	if err != nil {
		panic(err)
	}
	fmt.Println("affect row", affect)

	affect, err = e.ID(3).Cols("level").Update(&model.User{Name: "smith", Level: 2})
	if err != nil {
		panic(err)
	}
	fmt.Println("affect row", affect)

	affect, err = e.ID(2).MustCols("age").Update(&model.User{})
	if err != nil {
		panic(err)
	}
	fmt.Println("affect row", affect)

	affect, err = e.ID(5).AllCols().Update(&model.User{})
	if err != nil {
		panic(err)
	}
	fmt.Println("affect row", affect)

}

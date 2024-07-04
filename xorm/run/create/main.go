package main

import (
	"fmt"
	"learn_xorm/common"
	"learn_xorm/model"
)

func main() {
	users := make([]*model.User, 0, 6)
	users = append(users,
		&model.User{Name: "zz", Salt: "salt", Age: 18, Level: 2, Passwd: "12345"},
		&model.User{Name: "lj", Salt: "salt", Age: 1, Level: 1, Passwd: "12345"},
		&model.User{Name: "sd", Salt: "salt", Age: 8, Level: 6, Passwd: "12345"},
		&model.User{Name: "fs", Salt: "salt", Age: 2, Level: 5, Passwd: "12345"},
		&model.User{Name: "qw", Salt: "salt", Age: 8, Level: 4, Passwd: "12345"},
		&model.User{Name: "ee", Salt: "salt", Age: 7, Level: 3, Passwd: "12345"},
	)

	e := common.Engine()
	affect, err := e.Insert(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(affect)
}

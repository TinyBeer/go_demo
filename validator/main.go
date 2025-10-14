package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// 需要进行参数校验的结构体
type User struct {
	Name string `validate:"min=6,max=10"`  // Name字段 最少6个字符 最多10个字符
	Age  int    `validate:"min=1,max=100"` // Age字段 最小为1 最大为100
}

func main() {
	user := User{
		Name: "tom",
		Age:  101,
	}
	validate := validator.New()
	err := validate.Struct(user)
	fmt.Println(err)
}

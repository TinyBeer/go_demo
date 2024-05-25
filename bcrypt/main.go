package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "yourpassword"
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("散列结果长度：", len(bs))
	fmt.Println("散列结果", string(bs))

	cost, err := bcrypt.Cost(bs)
	if err != nil {
		panic("散列失败: err:" + err.Error())
	}
	fmt.Println("Hash Cost:", cost)

	err = bcrypt.CompareHashAndPassword(bs, []byte("12345sdfsdfsfsfsf6"))
	if err != bcrypt.ErrHashTooShort {
		fmt.Println(err)
	}

	err = bcrypt.CompareHashAndPassword(bs, []byte("yourpassword"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("密码匹配成功")
}

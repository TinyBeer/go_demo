package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	password := "123"
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bs))
	fmt.Println(len(bs))

	err = bcrypt.CompareHashAndPassword(bs, []byte("12345sdfsdfsfsfsf6"))
	if err != bcrypt.ErrHashTooShort {
		fmt.Println(err)
	}
}

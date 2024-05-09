package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 用于签名的字符串
var mySigningKey = []byte("tinybeer")

// GenRegisteredClaims 使用默认声明创建jwt
// 由于签发时间差异，所以不会出现相同的token
func GenRegisteredClaims() (string, error) {
	// 创建 Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 过期时间
		Issuer:    "beer",                                             // 签发人
	}
	// 生成token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	return token.SignedString(mySigningKey)
}

// ParseRegisteredClaims 解析jwt
func ValidateRegisteredClaims(tokenString string) bool {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil { // 解析token失败
		return false
	}
	return token.Valid
}

func main() {
	token, _ := GenRegisteredClaims()
	if ValidateRegisteredClaims(token) {
		fmt.Println("ok", token)
	}
}

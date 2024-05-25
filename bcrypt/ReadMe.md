bcrypt 包实现 Provos 和 Mazières 的 bcrypt 自适应散列算法。通常用于后端服务中密码脱敏存储。

# GenerateFromPassword

函数签名 `func GenerateFromPassword(password []byte, cost int) ([]byte, error)`
GenerateFromPassword 以给定的代价返回密码的 bcrypt 散列。const 范围`[4-31]`如果给定的 cost 小于 MinCost，则该 cost 将被设置为 DefaultCost`[10]`。如果 const 大于 MaxCost 则会报错。GenerateFromPassword 不接受长度超过 72 字节的密码，这是 bcrypt 操作的最长密码。

# Cost

函数签名`func Cost(hashedPassword []byte) (int, error)`
Cost 返回用于创建给定散列密码的散列成本。将来，当密码系统的哈希成本需要增加以适应更大的计算能力时，这个功能允许人们确定需要更新哪些密码。

# CompareHashAndPassword

函数签名 `func CompareHashAndPassword(hashedPassword, password []byte) error`
CompareHashAndPassword 将 bcrypt 散列密码与可能的明文密码进行比较。成功时返回 nil，失败时返回错误。

# 示例代码

```go
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

```

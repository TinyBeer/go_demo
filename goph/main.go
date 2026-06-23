package main

import (
	"log"

	"github.com/melbahja/goph"
)

func main() {

	// 创建 ssh 客户端i
	client, err := goph.New(
		// ssh 用户名
		"root",
		// ssh 服务器地址
		"10.160.162.42",
		// 使用密码方式鉴权
		goph.Password("dev123"))

	if err != nil {
		log.Fatal(err)
	}
	// 使用defer 关闭网络连接
	defer client.Close()

	// 执行 ls 命令
	out, err := client.Run("ls")
	if err != nil {
		log.Fatal(err)
	}

	// 打印执行结果
	print(string(out))
}

package main

import (
	"flag"
	"fmt"
)

func main() {
	var user string
	var pwd string
	var host string
	var port int
	flag.StringVar(&user, "u", "", "用户名，默认为空")
	flag.StringVar(&pwd, "p", "", "密码，默认为空")
	flag.StringVar(&host, "h", "localhost", "主机，默认为localhost")
	flag.IntVar(&port, "P", 3306, "端口号，默认为3306")
	flag.Parse()

	fmt.Println(user, pwd, host, port)
}

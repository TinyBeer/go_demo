package main

import (
	"fmt"
	"log"
	"net/rpc"
)

const serverAddress = "127.0.0.1"

func client() {
	client, err := rpc.Dial("tcp", serverAddress+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// 同步调用
	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	// 异步调用
	quotient := new(Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // 等带结果 will be equal to divCall

	if replyCall.Error != nil {
		fmt.Println("async call err:", replyCall.Error)
		return
	}
	args = replyCall.Args.(*Args)
	fmt.Printf("Arith.Divide: %d/%d=%d|%d\n", args.A, args.B, quotient.Quo, quotient.Rem)
}

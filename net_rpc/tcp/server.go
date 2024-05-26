package main

import (
	"log"
	"net"
	"net/rpc"
)

func server() {
	arith := new(Arith)
	rpc.Register(arith)
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	rpc.Accept(l)
}

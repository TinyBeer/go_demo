package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type Node struct {
	Next *Node
}

func main() {

	var name string = "tom"
	var age int = 18

	spew.Dump(name, age, Hello)

	spew.Fdump(os.Stdout, struct{ Price float64 }{88.8})

	str := spew.Sdump(&name, &age)
	fmt.Println(str)

	head := &Node{}
	head.Next = &Node{}
	head.Next.Next = head

	spew.Dump(head)
}

func Hello(name string) {
	fmt.Println("hello", name)
}

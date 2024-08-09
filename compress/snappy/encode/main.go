package main

import (
	"fmt"
	"os"

	"github.com/golang/snappy"
)

func main() {
	bs, err := os.ReadFile("./snapply_license")
	if err != nil {
		panic(err)
	}

	got := snappy.Encode(nil, bs)
	fmt.Printf("origin size:%d compressed:%d\n", len(bs), len(got))

	err = os.WriteFile("./snapply_license_compressed", got, 0666)
	if err != nil {
		panic(err)
	}
}

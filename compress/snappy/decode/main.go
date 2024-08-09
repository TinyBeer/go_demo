package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/golang/snappy"
)

func main() {
	compressed, err := os.ReadFile("./snapply_license_compressed")
	if err != nil {
		panic(err)
	}
	decompressed, err := snappy.Decode(nil, compressed)
	if err != nil {
		panic(err)
	}

	origin, err := os.ReadFile("./snapply_license")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(origin, decompressed) {
		panic("data break")
	}

	fmt.Println("ok")
}

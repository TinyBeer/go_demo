package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	uuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return
	}
	fmt.Printf("UUID: %v\n", uuid)
}

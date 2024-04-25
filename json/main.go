package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name   string  `json:"name,omitempty"`
	Age    int     `json:"age,omitempty"`
	Gender string  `json:"gender,omitempty"`
	Score  float64 `json:"score,omitempty"`
}

func main() {
	stu := Student{
		Name:   "tom",
		Age:    22,
		Gender: "M",
		Score:  0,
	}
	js, err := json.Marshal(stu)
	if err != nil {
		fmt.Println()
		return
	}
	fmt.Println(string(js))
	var nStu Student
	err = json.Unmarshal(js, &nStu)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nStu)
}

package main

import (
	"fmt"
	"reflect"
)

// 反射int
func ReflectInt() {
	var num int = 100

	// 获取num的Type
	rTyp := reflect.TypeOf(num)
	fmt.Printf("rType=%v type=%T\n", rTyp, rTyp)

	// 获取num的Vale
	rVal := reflect.ValueOf(num)
	fmt.Printf("rVal=%v type=%T\n", rVal, rVal)

	// 将rVal转回int变量
	iVal := rVal.Interface()
	fmt.Printf("iVal =%v type=%T\n", iVal, iVal)
	num2 := iVal.(int)
	fmt.Println("num2=", num2)

	// 修改num的值
	reflect.ValueOf(&num).Elem().SetInt(20)
	fmt.Println("num=", num)
}

type Student struct {
	Name string
	Age  int
}

// 放射结构体
func ReflectStruct() {
	stu := Student{
		Name: "tom",
		Age:  20,
	}

	// 获取stu的Type
	rTyp := reflect.TypeOf(stu)
	fmt.Printf("rType=%v type=%T\n", rTyp, rTyp)

	// 获取stu的Vale
	rVal := reflect.ValueOf(stu)
	fmt.Printf("rVal=%v type=%T\n", rVal, rVal)

	// 获取变量的kind
	kind := rTyp.Kind()
	kind2 := rVal.Kind()
	fmt.Printf("kind=%v kind2=%v\n", kind, kind2)

	// 将rVal转回Student变量
	iVal := rVal.Interface()
	fmt.Printf("iVal =%v type=%T\n", iVal, iVal)
	stu2 := iVal.(Student)
	fmt.Println("stu2=", stu2)
}

type Monster struct {
	Name  string `json:"name"`
	Age   int    `json:"monster_age"`
	Score float32
	Sex   string
}

func (s Monster) Print() {
	fmt.Println("---start---")
	fmt.Println(s)
	fmt.Println("---stop---")
}

func (s Monster) GetSum(n1, n2 int) int {
	return n1 + n2
}

func (s *Monster) Set(name string, age int, score float32, sex string) {
	s.Name = name
	s.Age = age
	s.Score = score
	s.Sex = sex
}

func ReflectStruct02() {
	mon := Monster{
		Name:  "黄鼠狼精",
		Age:   400,
		Score: 30.8,
		Sex:   "",
	}
	rType := reflect.TypeOf(mon)
	rVal := reflect.ValueOf(mon)
	kind := rVal.Kind()
	if kind != reflect.Struct {
		fmt.Println("变量不是一个结构体")
		return
	}

	fieldCnt := rVal.NumField()
	fmt.Printf("struct has %d fields\n", fieldCnt)
	for i := 0; i < fieldCnt; i++ {
		fmt.Printf("Field %d: %v\n", i, rVal.Field(i))
		tagVal := rType.Field(i).Tag.Get("json")
		if tagVal != "" {
			fmt.Printf("Field %d tag = %v\n", i, tagVal)
		}
	}

	methodCnt := rType.NumMethod()
	fmt.Printf("struct has %d methods\n", methodCnt)
	rVal.Method(1).Call(nil)

	var params []reflect.Value
	params = append(params, reflect.ValueOf(10), reflect.ValueOf(8))
	fmt.Println(rVal.Method(0).Call(params)[0])

	rVal = reflect.ValueOf(&mon)
	rVal.Elem().Field(0).SetString("白象精")
	fmt.Println(mon)
}

func main() {

	ReflectInt()

	ReflectStruct()

	ReflectStruct02()
}

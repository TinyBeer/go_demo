package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// 创建一个文件集
	fset := token.NewFileSet()

	// 解析源代码
	// src, err := os.ReadFile("./hello/hello.go")
	// if err != nil {
	// 	panic(err)
	// }
	// file, err := parser.ParseFile(fset, "hello.go", src, parser.ParseComments)
	// if err != nil {
	// 	panic(err)
	// }

	file, err := parser.ParseFile(fset, "./hello/hello.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// 打印语法树
	ast.Print(fset, file)
}

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strconv"
)

func main() {
	fileName := "./source.go"
	fset := token.NewFileSet()
	path, err := filepath.Abs(fileName)
	if err != nil {
		panic(err)
	}
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Print(fset, f)
	// 遍历语法树
	ast.Inspect(f, func(n ast.Node) bool {
		f, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		hasContext := false
		// 判断参数中是否包含context.Context类型
		for _, v := range f.Type.Params.List {
			if expr, ok := v.Type.(*ast.SelectorExpr); ok {
				if ident, ok := expr.X.(*ast.Ident); ok {
					if ident.Name == "context" {
						hasContext = true
					}
				}
			}
		}
		// 为没有context参数的方法添加context参数
		if !hasContext {
			ctxField := &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("ctx"),
				},
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("context"),
					Sel: ast.NewIdent("Context"),
				},
			}
			list := []*ast.Field{
				ctxField,
			}
			f.Type.Params.List = append(list, f.Type.Params.List...)
		}
		return false
	})
	addImport(f)
	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fset, f)
	if err != nil {
		log.Fatal(err)
	}
	// 输出Go代码
	fmt.Println(buffer.String())

}

// addImport 引入context包
func addImport(file *ast.File) {
	// 是否已经import
	hasImported := false
	for _, imptSpec := range file.Imports {
		if imptSpec.Path.Value == strconv.Quote("context") {
			hasImported = true
		}
	}
	// 如果没有import context，则import 没有考虑没有import的情况
	if !hasImported {
		for _, decl := range file.Decls {
			switch imp := decl.(type) {
			case *ast.GenDecl:
				if imp.Tok == token.IMPORT {
					imp.Specs = append(imp.Specs, &ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote("context"),
						},
					})
					return
				}
			}
		}
	}
}

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fileName := "./const.go"
	fset := token.NewFileSet()
	path, err := filepath.Abs(fileName)
	if err != nil {
		panic(err)
	}
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// 检查go:build标签
	modelFile := false
	for i, cg := range f.Comments {
		if i == 0 {
			for _, comment := range cg.List {
				if strings.HasPrefix(comment.Text, `//go:build `) || strings.HasPrefix(comment.Text, `// +build `) {
					if strings.Contains(comment.Text, "conster") {
						modelFile = true
						break
					}
				}
			}
		}
	}
	if !modelFile {
		fmt.Println(fileName, "is not a model file")
		return
	}

	ast.Print(fset, f)

	// 遍历语法树
	visitor := new(ConstDeclVisitor)
	visitor.Package = "main"
	ast.Walk(visitor, f)

	// 使用模板输出
	t, err := template.ParseFiles("./const.tmpl")
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile("./const_gen.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = t.Execute(file, visitor)
	if err != nil {
		panic(err)
	}
}

// 自定义Visiter
type ConstDeclVisitor struct {
	Package     string
	EnumSetList []*EnumSet
}

// 枚举变量数据暂存
type EnumSet struct {
	Type    string
	Members map[int]string
}

// 遍历方法
func (v *ConstDeclVisitor) Visit(node ast.Node) ast.Visitor {
	switch genDecl := node.(type) {
	// 判断否为*ast.GenDecl类型
	case *ast.GenDecl:
		// 根据*ast.GenDecl解析所需要的数据并存储
		set := new(EnumSet)
		if genDecl.Doc != nil && len(genDecl.Doc.List) != 0 {
			tmp, _ := strings.CutPrefix(genDecl.Doc.List[0].Text, `//`)
			set.Type = strings.TrimSpace(tmp)
		}

		for _, spec := range genDecl.Specs {
			switch valueSpec := spec.(type) {
			case *ast.ValueSpec:
				if set.Members == nil {
					set.Members = map[int]string{}
				}
				set.Members[valueSpec.Names[0].Obj.Data.(int)] = valueSpec.Names[0].Name
			}
		}
		if len(set.Members) != 0 {
			v.EnumSetList = append(v.EnumSetList, set)
		}

	case *ast.FuncDecl:
	}
	return v
}

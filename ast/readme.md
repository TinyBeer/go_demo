
# 词法分析

Golang 使用`go/scanner`中的Scanner类型来实现词法分析。通过词法分析将代码字符流转化为一个个的Token标记，方便语法分析器进行解析。

## Token标记

Token包含三部分内容，位置、标记、内容：

	* 位置 Token起始位置在文件中的定位
	* 标记 用于区分不同Token的标记 有唯一的数值进行区别
	* 内容 Token对应的文本内容

 | Token   | 说明                                           |
 | :------ | :--------------------------------------------- |
 | EOF     | 1 文件末尾 end of file                         |
 | COMMENT | 2 注释                                         |
 | IDENT   | 4 标记符 identifier 变量名、方法名等自定义标记 |
 | INT     | 5 整数 integer                                 |
 | FLOAT   | 6 浮点数字 float                               |
 | ...     | ...                                            |

具体内容参照示例结果理解，详细映射关系则可参考`go/token/token.go`文件。

## 词法分析示例

* 文件结构

	```bash
	.
	├── hello
	│   └── hello.go
	└── main.go
	```

* hello.go

	```go
	package main

	/*
	multi line comment
	*/
	import "fmt"

	// main hello
	func main() {
		content := "hello"
		fmt.Println(content)
	}
	```

* main.go

	```go
	package main

	import (
		"fmt"
		"go/scanner"
		"go/token"
		"io"
		"os"
	)

	func main() {
		oFile, err := os.OpenFile("./hello/hello.go", os.O_RDONLY, 0777)
		if err != nil {
			panic(err)
		}
		src, err := io.ReadAll(oFile)
		if err != nil {
			panic(err)
		}
		fset := token.NewFileSet()
		file := fset.AddFile(oFile.Name(), fset.Base(), len(src))
		var s scanner.Scanner
		s.Init(file, src, nil, scanner.ScanComments)
		for {
			pos, tok, lit := s.Scan()
			fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
			if tok == token.EOF {
				break
			}
		}

	}
	```

* 输出结果

	```bash
	./hello/hello.go:1:1    package "package"
	./hello/hello.go:1:9    IDENT   "main"
	./hello/hello.go:1:13   ;       "\n"
	./hello/hello.go:3:1    COMMENT "/*\n multi line comment\n*/"
	./hello/hello.go:6:1    import  "import"
	./hello/hello.go:6:8    STRING  "\"fmt\""
	./hello/hello.go:6:13   ;       "\n"
	./hello/hello.go:8:1    COMMENT "// main hello"
	./hello/hello.go:9:1    func    "func"
	./hello/hello.go:9:6    IDENT   "main"
	./hello/hello.go:9:10   (       ""
	./hello/hello.go:9:11   )       ""
	./hello/hello.go:9:13   {       ""
	./hello/hello.go:10:2   IDENT   "content"
	./hello/hello.go:10:10  :=      ""
	./hello/hello.go:10:13  STRING  "\"hello\""
	./hello/hello.go:10:20  ;       "\n"
	./hello/hello.go:11:2   IDENT   "fmt"
	./hello/hello.go:11:5   .       ""
	./hello/hello.go:11:6   IDENT   "Println"
	./hello/hello.go:11:13  (       ""
	./hello/hello.go:11:14  IDENT   "content"
	./hello/hello.go:11:21  )       ""
	./hello/hello.go:11:22  ;       "\n"
	./hello/hello.go:12:1   }       ""
	./hello/hello.go:12:2   ;       "\n"
	./hello/hello.go:12:3   EOF     ""
	```

# 语法分析

Golang中使用`go/ast`包来进行语法分析。`go/ast`包提供了对Go语言语法树的抽象表示。语法分析器将源代码转换为语法树（抽象语法树），语法树上每一个节点都代表着一种代码结构。通常我们遍历语法树上的节点，获取所需要的信息。

## 语法分析示例

* 文件结构

	```bash
	├── hello
	│   └── hello.go
	└── main.go
	```

* hello.go

	```go
	package main

	/*
	multi line comment
	*/
	import "fmt"

	// main hello
	func main() {
		content := "hello"
		fmt.Println(content)
	}
	```

* main.go

	```go
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
	```
  `ParseFile`提供了两种传入原文文件的方式：

  	* 一种是直接传入文件字节流，并提供自定义的文件名，如果src不为nil时使用这一方式。
  	* 另一种则是传入文件路径，在src为nil时会根据文件路径读取数据。

* 运行结果

	```bash
	0  *ast.File {
	1  .  Package: ./hello/hello.go:1:1
	2  .  Name: *ast.Ident {
	3  .  .  NamePos: ./hello/hello.go:1:9
	4  .  .  Name: "main"
	5  .  }
	6  .  Decls: []ast.Decl (len = 2) {
	7  .  .  0: *ast.GenDecl {
	8  .  .  .  Doc: *ast.CommentGroup {
	9  .  .  .  .  List: []*ast.Comment (len = 1) {
	10  .  .  .  .  .  0: *ast.Comment {
	11  .  .  .  .  .  .  Slash: ./hello/hello.go:3:1
	12  .  .  .  .  .  .  Text: "/*\n multi line comment\n*/"
	13  .  .  .  .  .  }
	14  .  .  .  .  }
	15  .  .  .  }
	16  .  .  .  TokPos: ./hello/hello.go:6:1
	17  .  .  .  Tok: import
	18  .  .  .  Lparen: -
	19  .  .  .  Specs: []ast.Spec (len = 1) {
	20  .  .  .  .  0: *ast.ImportSpec {
	21  .  .  .  .  .  Path: *ast.BasicLit {
	22  .  .  .  .  .  .  ValuePos: ./hello/hello.go:6:8
	23  .  .  .  .  .  .  Kind: STRING
	24  .  .  .  .  .  .  Value: "\"fmt\""
	25  .  .  .  .  .  }
	26  .  .  .  .  .  EndPos: -
	27  .  .  .  .  }
	28  .  .  .  }
	29  .  .  .  Rparen: -
	30  .  .  }
	31  .  .  1: *ast.FuncDecl {
	32  .  .  .  Doc: *ast.CommentGroup {
	33  .  .  .  .  List: []*ast.Comment (len = 1) {
	34  .  .  .  .  .  0: *ast.Comment {
	35  .  .  .  .  .  .  Slash: ./hello/hello.go:8:1
	36  .  .  .  .  .  .  Text: "// main hello"
	37  .  .  .  .  .  }
	38  .  .  .  .  }
	39  .  .  .  }
	40  .  .  .  Name: *ast.Ident {
	41  .  .  .  .  NamePos: ./hello/hello.go:9:6
	42  .  .  .  .  Name: "main"
	43  .  .  .  .  Obj: *ast.Object {
	44  .  .  .  .  .  Kind: func
	45  .  .  .  .  .  Name: "main"
	46  .  .  .  .  .  Decl: *(obj @ 31)
	47  .  .  .  .  }
	48  .  .  .  }
	49  .  .  .  Type: *ast.FuncType {
	50  .  .  .  .  Func: ./hello/hello.go:9:1
	51  .  .  .  .  Params: *ast.FieldList {
	52  .  .  .  .  .  Opening: ./hello/hello.go:9:10
	53  .  .  .  .  .  Closing: ./hello/hello.go:9:11
	54  .  .  .  .  }
	55  .  .  .  }
	56  .  .  .  Body: *ast.BlockStmt {
	57  .  .  .  .  Lbrace: ./hello/hello.go:9:13
	58  .  .  .  .  List: []ast.Stmt (len = 2) {
	59  .  .  .  .  .  0: *ast.AssignStmt {
	60  .  .  .  .  .  .  Lhs: []ast.Expr (len = 1) {
	61  .  .  .  .  .  .  .  0: *ast.Ident {
	62  .  .  .  .  .  .  .  .  NamePos: ./hello/hello.go:10:2
	63  .  .  .  .  .  .  .  .  Name: "content"
	64  .  .  .  .  .  .  .  .  Obj: *ast.Object {
	65  .  .  .  .  .  .  .  .  .  Kind: var
	66  .  .  .  .  .  .  .  .  .  Name: "content"
	67  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 59)
	68  .  .  .  .  .  .  .  .  }
	69  .  .  .  .  .  .  .  }
	70  .  .  .  .  .  .  }
	71  .  .  .  .  .  .  TokPos: ./hello/hello.go:10:10
	72  .  .  .  .  .  .  Tok: :=
	73  .  .  .  .  .  .  Rhs: []ast.Expr (len = 1) {
	74  .  .  .  .  .  .  .  0: *ast.BasicLit {
	75  .  .  .  .  .  .  .  .  ValuePos: ./hello/hello.go:10:13
	76  .  .  .  .  .  .  .  .  Kind: STRING
	77  .  .  .  .  .  .  .  .  Value: "\"hello\""
	78  .  .  .  .  .  .  .  }
	79  .  .  .  .  .  .  }
	80  .  .  .  .  .  }
	81  .  .  .  .  .  1: *ast.ExprStmt {
	82  .  .  .  .  .  .  X: *ast.CallExpr {
	83  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
	84  .  .  .  .  .  .  .  .  X: *ast.Ident {
	85  .  .  .  .  .  .  .  .  .  NamePos: ./hello/hello.go:11:2
	86  .  .  .  .  .  .  .  .  .  Name: "fmt"
	87  .  .  .  .  .  .  .  .  }
	88  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
	89  .  .  .  .  .  .  .  .  .  NamePos: ./hello/hello.go:11:6
	90  .  .  .  .  .  .  .  .  .  Name: "Println"
	91  .  .  .  .  .  .  .  .  }
	92  .  .  .  .  .  .  .  }
	93  .  .  .  .  .  .  .  Lparen: ./hello/hello.go:11:13
	94  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
	95  .  .  .  .  .  .  .  .  0: *ast.Ident {
	96  .  .  .  .  .  .  .  .  .  NamePos: ./hello/hello.go:11:14
	97  .  .  .  .  .  .  .  .  .  Name: "content"
	98  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 64)
	99  .  .  .  .  .  .  .  .  }
	100  .  .  .  .  .  .  .  }
	101  .  .  .  .  .  .  .  Ellipsis: -
	102  .  .  .  .  .  .  .  Rparen: ./hello/hello.go:11:21
	103  .  .  .  .  .  .  }
	104  .  .  .  .  .  }
	105  .  .  .  .  }
	106  .  .  .  .  Rbrace: ./hello/hello.go:12:1
	107  .  .  .  }
	108  .  .  }
	109  .  }
	110  .  FileStart: ./hello/hello.go:1:1
	111  .  FileEnd: ./hello/hello.go:12:3
	112  .  Scope: *ast.Scope {
	113  .  .  Objects: map[string]*ast.Object (len = 1) {
	114  .  .  .  "main": *(obj @ 43)
	115  .  .  }
	116  .  }
	117  .  Imports: []*ast.ImportSpec (len = 1) {
	118  .  .  0: *(obj @ 20)
	119  .  }
	120  .  Unresolved: []*ast.Ident (len = 1) {
	121  .  .  0: *(obj @ 84)
	122  .  }
	123  .  Comments: []*ast.CommentGroup (len = 2) {
	124  .  .  0: *(obj @ 8)
	125  .  .  1: *(obj @ 32)
	126  .  }
	127  .  GoVersion: ""
	128  }
	```
## 说明

`ParseFile`方法的返回值`ast.File`就是一棵语法树，每一个`ast.File`就是代表一个`go`文件的解析结果，其结构如下
```go
type File struct {
	Doc     *CommentGroup // associated documentation; or nil
	Package token.Pos     // position of "package" keyword
	Name    *Ident        // package name
	Decls   []Decl        // top-level declarations; or nil

	FileStart, FileEnd token.Pos       // start and end of entire file
	Scope              *Scope          // package scope (this file only). Deprecated: see Object
	Imports            []*ImportSpec   // imports in this file
	Unresolved         []*Ident        // unresolved identifiers in this file. Deprecated: see Object
	Comments           []*CommentGroup // list of all comments in the source file
	GoVersion          string          // minimum Go version required by //go:build or // +build directives
}
```

* `Doc`：顶部包说明注释，生成doc文件时会使用到。
* `Package`：`package`关键词，记录了位置信息。
* `Name`：类型为`*ast.Ident`，记录了包名所在位置和包名信息。
* `FileStart`, `FileEnd`：文件起始结束位置。
* `Imports`：类型为`[]*ImportSpec`，记录了所有引包信息
* `Comments`：类型为`[]*CommentGroup`，记录了所有的注释信息。
* `GoVersion`：类型为`string`，记录`Go:build`或`+build`指令要求的最小Go版本，仅在使用了这些指令之后才会有内容。
* `Decls`: 类型为`[]Decl`，记录了所有顶级声明结点的位置及相关信息。其中的`Decl`是一个接口所有声明类型的结点(包括`BadDecl`、`GenDecl`、`FuncDecl`)都实现了这个接口。
	* `BadDecl`：表示非法的声明, 仅记录其起始结束位置。
	```go
	BadDecl struct {
		From, To token.Pos // position range of bad declaration
	}
	```
	* `GenDecl`：代表包引用、常量变量以及类型声明。其中`Specs`代表具体的每一条声明语句(不包含括号部分)，包括三种类型`ImportSpec`、`ValueSpec`、`TypeSpec`
	* 。
	```go
	...
	GenDecl struct {
		Doc    *CommentGroup // associated documentation; or nil
		TokPos token.Pos     // position of Tok
		Tok    token.Token   // IMPORT, CONST, TYPE, or VAR
		Lparen token.Pos     // position of '(', if any
		Specs  []Spec
		Rparen token.Pos // position of ')', if any
	}
	...
	```
	* `ImportSpec`: 包引用声明结点
		```go
		...
		ImportSpec struct {
			Doc     *CommentGroup // associated documentation; or nil
			Name    *Ident        // local package name (including "."); or nil
			Path    *BasicLit     // import path
			Comment *CommentGroup // line comments; or nil
			EndPos  token.Pos     // end of spec (overrides Path.Pos if nonzero)
		}
		...
		```
	* `ValueSpec`: 常量变量声明结点
		```go
		...
		ValueSpec struct {
			Doc     *CommentGroup // associated documentation; or nil
			Names   []*Ident      // value names (len(Names) > 0)
			Type    Expr          // value type; or nil
			Values  []Expr        // initial values; or nil
			Comment *CommentGroup // line comments; or nil
		}
		...
		```
	* `TypeSpec`： 类型声明结点
		```go
		...
		TypeSpec struct {
			Doc        *CommentGroup // associated documentation; or nil
			Name       *Ident        // type name
			TypeParams *FieldList    // type parameters; or nil
			Assign     token.Pos     // position of '=', if any
			Type       Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
			Comment    *CommentGroup // line comments; or nil
		}
		...
		```

	* `FuncDecl`: 代表函数或者方法的声明。
	```go
	...
	FuncDecl struct {
		Doc  *CommentGroup // associated documentation; or nil
		Recv *FieldList    // receiver (methods); or nil (functions)
		Name *Ident        // function/method name
		Type *FuncType     // function signature: type and value parameters, results, and position of "func" keyword
		Body *BlockStmt    // function body; or nil for external (non-Go) function
	}
	...
	```
除了这些结点以外，还有两种基础结点`Expr`、`Stmt`参与构成语法树，分别代表表达式和语句。他们的种类很多（参看下方代码），几乎每一种golang语法都有其对应的`Expr`或`Stmt`结点(有一些是多种语法对应一种结点)，这里就不进行一一说明了。使用时根据需求翻阅源码即可。

```go
// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.
type (
	// A BadExpr node is a placeholder for an expression containing
	// syntax errors for which a correct expression node cannot be
	// created.
	//
	BadExpr struct {
		From, To token.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
		Obj     *Object   // denoted object, or nil. Deprecated: see Object.
	}

	// An Ellipsis node stands for the "..." type in a
	// parameter list or the "..." length in an array type.
	//
	Ellipsis struct {
		Ellipsis token.Pos // position of "..."
		Elt      Expr      // ellipsis element type (parameter lists only); or nil
	}

....

// ----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	//
	BadStmt struct {
		From, To token.Pos // position range of bad statement
	}

	// A DeclStmt node represents a declaration in a statement list.
	DeclStmt struct {
		Decl Decl // *GenDecl with CONST, TYPE, or VAR token
	}

	// An EmptyStmt node represents an empty statement.
	// The "position" of the empty statement is the position
	// of the immediately following (explicit or implicit) semicolon.
	//
	EmptyStmt struct {
		Semicolon token.Pos // position of following ";"
		Implicit  bool      // if set, ";" was omitted in the source
	}

	// A LabeledStmt node represents a labeled statement.
	LabeledStmt struct {
		Label *Ident
		Colon token.Pos // position of ":"
		Stmt  Stmt
	}
...
```

# 如何使用 AST

对于语法树的访问，golang提供了`ast.Inspect(node ast.Node, f func(ast.Node) bool)`和`ast.Walk(v ast.Visitor, node ast.Node)`方法进行遍历语法树上的结点。其中`Inspect`是通过`Walk`方法实现的。我们可以通过接口变量断言的方式找到自己需要的结点，从而进行相应的处理。下面会通过两种场景演示两种不同的玩法。

## 生产新文件

将go语言作为描述语言，生成基于go语言定义的代码文件。
文件 const.go 用于描述要生成的代码。这里希望将`const`域内的常量定义为指定类型的变量，变量类型根据注释定义。同时为变量实现`String`方法。
这里还用到了`go:build`标签语法，用来区分哪些描述文件，(当然也可通过修改文件后缀实现)，同时也可以避免当成go代码文件被解析。
```go
//go:build conster
// +build conster

package main

// PlanType
const (
	PlanType_Daily = iota
	PlanType_Weekly
	PlanType_Monthly
)

// TodoType
const (
	TodoType_Times = iota
	TodoType_Duration
	TodoType_TimesAndDuration
)

```

main.go
```go
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

	// 使用模板输出  当然，使用其他方式拼接内容也可以实现相同功能
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

```
const.tmpl 模板文件
```tmpl
// Code generated by conster. DO NOT EDIT.

package {{ .Package }}
{{ range $EnumSet := .EnumSetList }}
type {{ $EnumSet.Type }} int 
var ( {{ range $key, $value := $EnumSet.Members }}
    {{$value}} {{$EnumSet.Type}}  = {{$key}}{{ end }}
)
func (v {{ .Type }})String() string {
    switch v { {{ range $key, $value := .Members }}
    case {{$key}}: 
      return "{{$value}}"{{ end }}
    default:
      return ""
    }
}
{{ end }}
```

const_gen.go 输出结果
```go
// Code generated by conster. DO NOT EDIT.

package main

type PlanType int 
var ( 
    PlanType_Daily PlanType  = 0
    PlanType_Weekly PlanType  = 1
    PlanType_Monthly PlanType  = 2
)
func (v PlanType)String() string {
    switch v { 
    case 0: 
      return "PlanType_Daily"
    case 1: 
      return "PlanType_Weekly"
    case 2: 
      return "PlanType_Monthly"
    default:
      return ""
    }
}

type TodoType int 
var ( 
    TodoType_Times TodoType  = 0
    TodoType_Duration TodoType  = 1
    TodoType_TimesAndDuration TodoType  = 2
)
func (v TodoType)String() string {
    switch v { 
    case 0: 
      return "TodoType_Times"
    case 1: 
      return "TodoType_Duration"
    case 2: 
      return "TodoType_TimesAndDuration"
    default:
      return ""
    }
}

```

## 修改语法树

我们也可以直接修改语法树，从而实现对源代码的修改。

场景： 为没有添加`context.Contex`的方法、函数添加上下问参数。

source.go 
```go
package main

type Foo struct {
}

func (*Foo) NeedContext() {

}

func ContextWanted(name string) string {
	return name
}

```

main.go
```go
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

```

output 
```go
package main

import (
        "context"
        "strings"
)

type Foo struct {
}

func (*Foo) NeedContext(ctx context.Context) {

}

func (*Foo) NotNeedContext(ctx context.Context) {

}

func ContextWanted(ctx context.Context, name string) string {
        return strings.TrimSpace(name)
}
```
可以看到我们所有方法和函数的第一个参数都变成了`context.Context`。
补充说明：虽然直接修改语法树可以方便的对源码进行修改，但是这种方法也存在一些问题。如注释位置变化会发生变化。这是由于修改了文件后，节点的起始终止位置发生了变化。对于非注释节点，语法树能够正确的调整他们的位置，但却不能自动调整注释节点的位置。如果我们想要让注释出现在正确的位置上，我们必须手动设置节点`Pos`和`End`。此外，偶尔会出现入参后面多出一个逗号的情况。

# 参考

[ast 源码文档](https://pkg.go.dev/go/ast)
[Golang AST语法树使用教程及示例](https://juejin.cn/post/6844903982683389960)

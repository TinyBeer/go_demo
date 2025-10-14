# 简介
`validator` 是一款使用广泛的参数校验库，主要功能是通过结构体 `Tag` 配置校验规则，实现对结构体参数的校验。

# 快速开始
## 安装
```bash
go get github.com/go-playground/validator/v10
```

## 引用
```golang
import "github.com/go-playground/validator/v10"
```

## 基础使用
```golang
package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// 需要进行参数校验的结构体
type User struct {
	Name string `validate:"min=6,max=10"`  // Name字段 最少6个字符 最多10个字符
	Age  int    `validate:"min=1,max=100"` // Age字段 最小为1 最大为100
}

func main() {
	user := User{
		Name: "tom",
		Age:  101,
	}
	validate := validator.New()
	err := validate.Struct(user)
	fmt.Println(err)
}
```
```bash
# 执行结果
$ go run main.go
Key: 'User.Name' Error:Field validation for 'Name' failed on the 'min' tag
Key: 'User.Age' Error:Field validation for 'Age' failed on the 'max' tag
```

# 常用
## 常用约束项

|Tag|描述|
|:----|:----|
|len|长度为|
|min/max|最小/大值|
|eq(eq_ignore_case)/neq(ne_ignore_case)|等于(忽略大小写)/不等于(忽略大小写)|
|gt/lt(gte/lte)|大于/小于(大于等于/小于等于)|
|required|不为零值|
|oneof|枚举值之一|
|eqfield|等于同级字段|
|contains/excludes|包含/不包含子串|
|startswith/endswith|以子串开始/结束|
|unqiue|唯一性约束|
|email|邮箱格式|

[更多](https://github.com/go-playground/validator?tab=readme-ov-file#fields)

## 使用技巧
```golang
type Peopole struct {
	Name      string   `validate:"min=3,max=2"`       // 长度范围
	Gender    string   `validate:"oneof=male female"` // 性别
	Email     string   `validate:"required,email"`    //邮件
	Password  string   `validate:"min=6"`
	Password2 string   `validate:"eqfield=Password"` //二次密码
	Hobbies   []string `validate:"unique"`           // 唯一性
}

```
# 自定义约束
<!-- todo -->

# 错误处理
<!-- todo -->


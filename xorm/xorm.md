[xorm](https://xorm.io/xorm/) 是一个简单但强大的 Go 语言 ORM 库。
# 快速开始
```sh
# 安装xorm库
go get xorm.io/xorm
# 根据操作的数据库下载相应驱动
# 下载操作mysql的驱动
go get github.com/go-sql-driver/mysql
```

# 创建引擎
所有操作均在 `ORM引擎` 上进行，所以需要先创建并配置好ORM引擎。XORM 支持 `Engine` 和 `Engine Group` 两种引擎。 
* `Engine` 引擎 用于对单个数据库进行操作。
* `Engine Group` 引擎用于对读写分离的数据库或者负载均衡的数据库进行操作。

`Engine` 引擎和 `Engine Group` 引擎的API基本相同，所有适用于 `Engine` 的 API 基本上都适用于 `Engine Group`，并且可以比较容易的从 `Engine` 迁移到 `Engine Group` 。
### Engine
Engine 的参数与`sql.Open`参数相同，同样的也需要引入对应数据库的驱动包，相关驱动包可以在[官网](https://xorm.io)查看`Databases/Drivers supported`。
```go
// 没有引入会报错：sql: unknown driver "mysql" (forgotten import?)
import	_ "github.com/go-sql-driver/mysql"
...

engine, err := xorm.NewEngine("mysql", "root:123456@/test?charset=utf8mb4")
...
```

你也可以用 `NewEngineWithParams`, `NewEngineWithDB` 和 `NewEngineWithDialectAndDB` 来创建引擎。
一般情况下如果只操作一个数据库，只需要创建一个 engine 即可。engine 是 GoRoutine 安全的。
创建完成 engine 之后，并没有立即连接数据库，此时可以通过 engine.Ping() 或者 engine.PingContext() 来进行数据库的连接测试是否可以连接到数据库。另外对于某些数据库有连接超时设置的，可以通过起一个定期Ping的Go程来保持连接鲜活。
### Engine Group
<!-- TODO: 补充内容 -->

### 日志
日志是一个接口，通过设置日志，可以显示SQL，警告以及错误等，默认的显示级别为 INFO。
* engine.ShowSQL(true)，则会在控制台打印出生成的SQL语句；
* engine.Logger().SetLevel(log.LOG_DEBUG)，则会在控制台打印调试及以上的信息；
* engine.SetLogger(log.NewSimpleLogger(f))，修改日志输出目标，`NewSimpleLogger(w io.Writer)`接收一个io.Writer接口来将数据写入到对应的设施中。
  
### 连接池
engine内部支持连接池接口和对应的函数。
* 如果需要设置连接池的空闲数大小，可以使用 engine.SetMaxIdleConns() 来实现。
* 如果需要设置最大打开连接数，则可以使用 engine.SetMaxOpenConns() 来实现。
* 如果需要设置连接的最大生存时间，则可以使用 engine.SetConnMaxLifetime() 来实现。

# 表操作
## 表定义
xorm使用结构体进行表结构定义。
表名的设置主要有两种方式：
* Mapper映射规则，在没有`TableName`方法时启用。xorm内置了三种Mapper规则。
  * SnakeMapper （默认） 支持struct为驼峰式命名，表结构为下划线命名之间的转换，这个是默认的Maper；
  * SameMapper 支持结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名；
  * GonicMapper 和SnakeMapper很类似，但是对于特定词支持更好，比如ID会翻译成id而不是i_d。
  可以通过`engine.SetMapper(names.GonicMapper{})`方式设置，此外也可以自己实现Mapper以满足特殊的需求。
  默认情况下，映射规则在表名和字段名中共用，如果需要区别开可以使用以下方式：
  ```go
  engine.SetTableMapper(names.SameMapper{})
  engine.SetColumnMapper(names.SnakeMapper{})
  ```
  此外通过 `names.NewPrefixMapper(names.SnakeMapper{}, "prefix")` 可以创建一个在 Mapper 的基础上在命名中添加统一的前缀。
* 结构体实现`TableName() string`方法，优先使用。

字段Tag规则:
1. 字段名 
    字段的名称，可选，如不写，则自动根据field名字和转换规则命名。如果与关键字冲突需要使用单引号包裹，否则之际使用即可。
    ```go
    Name string  `xorm:"varchar(25) notnull unique 'usr_name' comment('姓名')"`
    ```
2. 关键字

| 关键字                                                                  | 说明                                                                                               |
| :---------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------- |
| pk                                                                      | 是否是主键，如果有多个字段都使用了此标记，则为复合主键                                             |
| 当前支持[30多种字段类型](https://xorm.io/zh/docs/chapter-02/4.columns/) | 字段类型                                                                                           |
| autoincr                                                                | 是否是自增                                                                                         |
| [ not ]null 或 notnull                                                  | 是否可以为空                                                                                       |
| unique或unique(uniquename)                                              | 括号中为联合唯一索引的名字，如果多个字段uniquename相同，则这些uniquename相同的字段组成联合唯一索引 |
| index或index(indexname)                                                 | 括号中为联合索引的名字，如果多个字段indexname相同，则这些indexname相同的字段组成联合索引           |
| extends                                                                 | 应用于员结构体之上，表示此结构体的所有成员也映射到数据库中，extends可加载无限级                    |
| -                                                                       | 这个Field将不进行字段映射                                                                          |
| ->                                                                      | 这个Field将只写入到数据库而不从数据库读取                                                          |
| <-                                                                      | 这个Field将只从数据库读取，而不写入到数据库                                                        |
| created                                                                 | 这个Field将在Insert时自动赋值为当前时间                                                            |
| updated                                                                 | 这个Field将在Insert或Update时自动赋值为当前时间                                                    |
| deleted                                                                 | 这个Field将在Delete时设置为当前时间，并且当前记录不删除                                            |
| version                                                                 | 这个Field将会在insert时默认为1，每次更新自动加1                                                    |
| default                                                                 | 0或default(0)	设置默认值，紧跟的内容如果是Varchar等需要加上单引号                                  |
| json                                                                    | 表示内容将先转成Json格式，然后存储到数据库中，数据库中的字段类型可以为Text或者二进制               |
| comment                                                                 | 设置字段的注释（当前仅支持mysql）                                                                  |

> 单主键当前支持int32,int,int64,uint32,uint,uint64,string这7种Go的数据类型，复合主键支持这7种Go的数据类型的组合。

## 建表
建表可以使用`Sync`方法, `Sync2` 已弃用。

> Sync 不会对已有表结构的字段进行任何的更改、删除操作，而只会进行增加字段操作。 

```go
err := engine.Sync(new(User), new(Group))
```

> 老版本如果发现有字符删减，会打印Warn级别日志进行提示。新版本则不会，如果有需求，可以使用`SyncWithOptions`方法开启警告。 

```go
// 开启警告
err = engine.SyncWithOptions(xorm.SyncOptions{
		WarnIfDatabaseColumnMissed: true,
		IgnoreConstrains:           false,
		IgnoreIndices:              false,
		IgnoreDropIndices:          false,
	}, new(User), new(Group))
```

> 新版本Sync方法会一次性拉取所有数据库的表数据，与注入结构体对比。推荐将所有表结构一次性注入Sync方法，以提高同步效率。

```go
// 不推荐
engine.Sync(new(User))
engine.Sync(new(Group))

// 推荐
engine.Sync(new(User), new(Group))

```

# 增删改查
## 插入数据
使用`engine.Insert()`方法，可以插入单条数据，也可以批量插入多条数据：
```go
// 单挑插入，会更新user中的id字段
user := &User{Name: "zz", Salt: "salt", Age: 18, Level: 2, Passwd: "12345"}
affect, err := engine.Insert(user)

// 批量插入，不会修改id字段
users := make([]*User, 0, 6)
users = append(users,
  &User{Name: "zz", Salt: "salt", Age: 18, Level: 2, Passwd: "12345"},
  &User{Name: "lj", Salt: "salt", Age: 1, Level: 1, Passwd: "12345"},
  &User{Name: "sd", Salt: "salt", Age: 8, Level: 6, Passwd: "12345"},
  &User{Name: "fs", Salt: "salt", Age: 2, Level: 5, Passwd: "12345"},
  &User{Name: "qw", Salt: "salt", Age: 8, Level: 4, Passwd: "12345"},
  &User{Name: "ee", Salt: "salt", Age: 7, Level: 3, Passwd: "12345"},
)
affect, err = engine.Insert(&users)
```
> Insert传入切片时，数据分成多条SQL进行插入，如果其中一条插入时出错，可能造成后续数据不能继续插入。
>
> Insert传入切片指针时，数据会拼接为一条SQl执行能够插入，由于各数据库SQL长度限制，需要注意不要使用太长的切片内容。如果数据量太大，需要分片进行插入。官方建议每条长度不要超过150。
> 
> Insert支持多个不同类型的参数。可以进行多张表的插入。

## 查询
查询主要使用`Get`和`Find`方法，`Get`查找一个结果，`Find`则为多个结果，结果会保存的出入参数中。
```go
// 第一个返回值为是否查询到结果
ok, err := e.ID(2).Get(u)

err = e.Find(&us)
```
查询参数
* `Where("a = ? AND b = ?", 1, 2)`: where 条件  
* `Alias("o").Where("o.name = ?", name)`: 表别名
* `.Where(...).And(...)`: AND
* `.Asc("id")` `.Desc("time")` `OrderBy(string)`：结果排序，可以组合使用
* `.ID(1)`：查询主键
* `Or(interface{}, …interface{})`： OR条件
* `Select(string)`：Select内容
* `SQL(string, …interface{})`：SQL内容
* `.In("cloumn", 1, 2, 3)`: IN
* `.Cols("age", "name")`： 查询或更新的字段
* `.AllCols()`：查询或更新所有字段
* `MustCols(…string)`：某些字段必须更新，一般与Update配合使用。
* `.Omit("age", "gender")`: 和cols相反，此函数指定排除某些指定的字段。注意：此方法和Cols方法不可同时使用。
* `.Distinct("age", "department")`: 归类去重
* `Table(nameOrStructPtr interface{})`：传入表名称或者结构体指针，如果传入的是结构体指针，则按照IMapper的规则提取出表名
* `Limit(int, …int)`：限制获取的数目，第一个参数为条数，第二个参数表示开始位置，如果不传则为0
* `Top(int)`：相当于Limit(int, 0)
* `Join(string,interface{},string)`
* `GroupBy(string)`
* `Having(string)`
其他方法
* `Count`: 统计数据数量
* `Exist`： 判断某个记录是否存在，比Get，Exist性能更好。
* `Iterate`：Iterate方法提供逐条执行查询到的记录的方法
* `Rows`： 同样供逐条执行查询到的记录的方法，不过Rows更加灵活好用。
* `Sum`：求和数据可以使用`Sum`, `SumInt`, `Sums` 和 `SumsInt` 四个方法
## 删除
```go
affected, err := engine.Where("name = ?", "lzy").Delete(&User{})
```
> `engine.Delete(new(User))` 不会进行删除操作，需要使用`engine.Where("1=1").Delete(new(User))`
> 
> 使用软删除，需要在xorm标记中使用deleted标记 

```go
DeletedAt time.Time `xorm:"deleted"`
```

## 修改

更新通过`engine.Update()`实现，可以传入结构指针或`map[string]interface{}`。对于传入结构体指针的情况，xorm只会更新非空的字段。如果一定要更新空字段，需要使用`Cols()`方法显示指定更新的列。使用`Cols()`方法指定列后，即使字段为空也会更新


```go
affect, err := engine.ID(1).Update(&model.User{Name: "jack"})
```

```go
affect, err := engine.Table(&model.User{}).ID(2).Update(map[string]interface{}{"age": 16})
```

```go
affect, err := engine.ID(3).Cols("level").Update(&model.User{Name: "smith", Level: 2})
```
# 执行原始SQL
## 查询
```go
sql := "select * from userinfo"
results, err := engine.Query(sql)
```
返回值 `results` 为 `[]map[string][]byte`
`Query` 的参数也允许传入 `*builder.Buidler` 对象

此外还可以使用`QueryInterface`（返回类型`[]map[string]interface{}`） 和 `QueryString`（返回类型`[]map[string]string`）

## 执行命令
```go
sql = "update `userinfo` set username=? where id=?"
res, err := engine.Exec(sql, "xiaolun", 1)
```
# 事务

当使用事务处理时，需要创建 Session 对象。在进行事务处理时，可以混用 ORM 方法和 SQL 方法。
```go
session := engine.NewSession()
defer session.Close()

// add Begin() before any action
if err := session.Begin(); err != nil {
    return err
}

user1 := Userinfo{Username: "xiaoxiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
if _, err := session.Insert(&user1); err != nil {
    return err
}
user2 := Userinfo{Username: "yyy"}
if _, err = session.Where("id = ?", 2).Update(&user2); err != nil {
    return err
}

if _, err = session.Exec("delete from userinfo where username = ?", user2.Username); err != nil {
    return err
}

// add Commit() after all actions
return session.Commit()
```

# reverse 工具
## 安装
```sh
go get xorm.io/reverse
go install xorm.io/reverse
```
## 使用默认模板
```sh
reverse -f example/custom.yml
```

```yaml
kind: reverse
name: test
source:
  database: mysql
  conn_str: 'root:123456@/test?charset=utf8mb4'
targets:
- type: codes
  language: golang
  output_dir: ./models
```
生成的go文件 `./models/models.go`
```go
package models

import (
	"time"
)

type MyUser struct {
	ID      int64     `xorm:"not null pk autoincr BIGINT(20)"`
	Name    string    `xorm:"VARCHAR(255)"`
	Salt    string    `xorm:"VARCHAR(255)"`
	Age     int       `xorm:"INT(11)"`
	Level   int       `xorm:"INT(11)"`
	Passwd  string    `xorm:"VARCHAR(200)"`
	Created time.Time `xorm:"DATETIME"`
	Updated time.Time `xorm:"DATETIME"`
}
```
## 更复杂的配置
```yaml
# complex.yaml
kind: reverse
name: mydb
source:
  database: mysql
  conn_str: 'root:123456@/test?charset=utf8mb4'
targets:
- type: codes
  # include_tables: # 需要导出的表
  #   - a
  #   - b
  # exclude_tables: # 需要排除的表
  #   - c
  table_mapper: snake # 表名映射规则
  column_mapper: snake # 字段映射规则
  table_prefix: "my_" # 表名前缀  默认会去除掉的部分
  multiple_files: true # 生产多份文件
  language: golang
  template: | # 代码模板 使用的是golang模板语法
    package models

    {{$ilen := len .Imports}}
    {{if gt $ilen 0}}
    import (
      {{range .Imports}}"{{.}}"{{end}}
    )
    {{end}}

    {{range .Tables}}
    type {{TableMapper .Name}} struct {
    {{$table := .}}
    {{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{ColumnMapper $col.Name}}	{{Type $col}} `{{Tag $table $col}}`
    {{end}}
    }

    func (m *{{TableMapper .Name}}) TableName() string {
    	return "{{$table.Name}}"
    }
    {{end}}
  template_path: ./template/goxorm.tmpl # 模板文件路径  优先级低于template字段
  output_dir: ./models # 生产文件路径
```
生成的go文件 `./models/user.go`
```go
package models

import (
	"time"
)

type User struct {
	ID      int64     `xorm:"not null pk autoincr BIGINT(20)"`
	Name    string    `xorm:"VARCHAR(255)"`
	Salt    string    `xorm:"VARCHAR(255)"`
	Age     int       `xorm:"INT(11)"`
	Level   int       `xorm:"INT(11)"`
	Passwd  string    `xorm:"VARCHAR(200)"`
	Created time.Time `xorm:"DATETIME"`
	Updated time.Time `xorm:"DATETIME"`
}

func (m *User) TableName() string {
	return "user"
}
```

# 参考
[Go 每日一库之 xorm](https://darjun.github.io/2020/05/07/godailylib/xorm/)
[XORM - eXtra ORM for Go](https://xorm.io/zh/)


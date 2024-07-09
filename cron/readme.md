[cron](https://github.com/robfig/cron)一个用于管理定时任务的库，用 Go 实现 Linux 中crontab这个命令的效果。
# 快速开始

## 安装
```sh
go get github.com/robfig/cron/v3@v3.0.0
```

## 示例代码
```go
package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	start := time.Now()
	c.AddFunc("@every 1s", func() {
		fmt.Println(time.Since(start))
	})
	c.Start()

	time.Sleep(time.Second * 10)
}
``` 
使用者向cron对象中注册函数，cron会根据对于的执行计划启动协程执行他们。这里注册了一个打印时长的函数，计划为每间隔一秒执行一次。执行效果如下：
```sh
$ go run ./demo/quick_start
359.657307ms
1.359776034s
2.359991505s
3.359997859s
4.360006728s
5.360005174s
6.360014508s
7.36005255s
8.359200507s
9.359380883s
```
# 时间格式

## 时间字段
cron库支持用 6 个空格分隔的字段来表示时间。

> 默认情况下cron解释器使用后五个字段，如果需要开启对秒的支持，可以使用`cron.New(cron.WithSeconds())`创建cron对象。

```sh
# ┌───────────── 秒 (0-59)
# │ ┌───────────── 分 (0–59)
# │ │ ┌───────────── 时 (0–23)
# │ │ │ ┌───────────── 几号 (1–31)
# │ │ │ │ ┌───────────── 月份 (1–12)
# │ │ │ │ │ ┌───────────── 周几 (0–6) (Sunday to Saturday;
# │ │ │ │ │ │                                   
# │ │ │ │ │ │
# │ │ │ │ │ │
# *  *  *  *  *  * <时间表达式>

```
| 字段 | 值              | 支持的特殊字符 |
| ---- | --------------- | -------------- |
| 秒   | 0-59            | * / , -        |
| 分   | 0-59            | * / , -        |
| 时   | 0-23            | * / , -        |
| 几号 | 1-31            | * / , - ?      |
| 月份 | 1-12 or JAN-DEC | * / , -        |
| 周几 | 0-6 or SUN-SAT  | * / , - ?      |

* 月份和周几大小写不敏感(SUN,Sun,sun是等效的)
* 特殊符号
  * `*` 用于匹配所有值
  * `/` 用于指定步长，如`1-5/2`从1开始每间隔两个时间单位触发，到5为止
  * `,` 用于表示或关系
  * `-` 用于表示范围
  * `?` 在月份和周几中替代`*`表示任意一天

## 预定义计划
cron提供以下预定义计划：
| 格式                   | 说明                 | 等价于      |
| ---------------------- | -------------------- | ----------- |
| @yearly (或 @annually) | 每一年的一月一日零时 | 0 0 0 1 1 * |
| @monthly               | 每月一号零时         | 0 0 0 1 * * |
| @weekly                | 每周一零时           | 0 0 0 * * 0 |
| @daily (或 @midnight)  | 每天零时             | 0 0 0 * * * |
| @hourly                | 每小时零分零秒       | 0 0 * * * * |

## 间隔时间
`@every <duration>` 按照固定时间间隔执行。`duration`需要可以被[`time.ParseDuration`](http://golang.org/pkg/time/#ParseDuration)解析。

# 时区
cron默认情况下使用机器本地时区。同时也支持通过`WithLocation`配置全局时区以及在时间表达式前使用关键字(`TZ=`或者`CRON_TZ=`)指定时区。
```go
// 配置全局时区
loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
c := cron.New(cron.WithLocation(loc))

// 指定任务时区
c.AddFunc("CRON_TZ=Asia/Tokyo 13 11 * * *", func() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	fmt.Println("tokyo time", time.Now().In(loc))
})
```

# Job
cron提供使用Job接口注册计划任务的方法`AddJob(Job)`，相较于使用无参数的函数，可以更方便的管理变量。事实上`AddFunc`也是通过`AddJob`注册计划任务的(将无参数函数转化为实现了Job接口的FuncJob)。
```go
// Job is an interface for submitted cron jobs.
type Job interface {
	Run()
}
```
配合sync.Mutex进行并发控制。

> cron会创建新的go协程来执行任务，对于资源的并发访问控制，需要自行进行处理。

```go
func main() {
	c := cron.New()
	myjob := new(MyJob)
	c.AddJob("@every 1s", myjob)
	c.AddJob("@every 2s", myjob)
	c.Start()

	time.Sleep(time.Second * 10)
}

type MyJob struct {
	mu sync.Mutex
	id int
}

func (m *MyJob) Run() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.id++
	fmt.Println("work with lock", m.id, time.Now())
	time.Sleep(time.Second * 3)
	fmt.Println("finish", m.id, time.Now())
}
```

# 其他选项 
## WithParser
用于配置自定义解析器，cron提供了以下配置项目：
```go
const (
	Second         ParseOption = 1 << iota // Seconds field, default 0
	SecondOptional                         // Optional seconds field, default 0
	Minute                                 // Minutes field, default 0
	Hour                                   // Hours field, default 0
	Dom                                    // Day of month field, default *
	Month                                  // Month field, default *
	Dow                                    // Day of week field, default *
	DowOptional                            // Optional day of week field, default *
	Descriptor                             // Allow descriptors such as @monthly, @weekly, etc.
)
```
其中SecondOptional和DowOptional为秒字段和周几字段是否为可选，两者不能同时使用。
使用自定义解析器，我们可以只保留关注的字段，简化需要使用的时间表达式。

使用示例
```go
// 仅关注每天内的任务 并且支持 秒字段可选
c := cron.New(cron.WithParser(cron.NewParser(cron.Second | cron.SecondOptional | cron.Minute | cron.Hour)))
// 每分钟的第3秒执行
c.AddFunc("3 * *", func() {
	fmt.Println(time.Now())
})

// 每分钟第0秒执行
c.AddFunc("* *", func() {
	fmt.Println(time.Now())
})
```

## WithChain
Job装饰器，cron提供了三种内置的装饰器:
* Recover：捕获内部Job产生的 panic；
* DelayIfStillRunning：触发时，如果上一次任务还未执行完成，则等待上一次任务完成之后再执行；
* SkipIfStillRunning：触发时，如果上一次任务还未完成，则跳过此次执行。

全局配置：
```go
c := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
c.AddFunc("@every 1s", func() {
	panic(time.Now())
})
c.Start()
```

为任务单独配置：
```go
c.AddJob(
	"@every 1s",
	cron.NewChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	).Then(
		cron.FuncJob(func() {
			fmt.Println("hahha")
			time.Sleep(time.Second * 2)
		}),
	),
)
```

## WithLogger
用于配置cron的日志打印，默认情况下cron不会打印info级别的日志。
```go
// 打印详细的执行详细
c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
```
<!-- TODO # 源码 -->

# 参考
[cron document](https://pkg.go.dev/github.com/robfig/cron)
[cron wikipedia page](https://en.wikipedia.org/wiki/Cron)
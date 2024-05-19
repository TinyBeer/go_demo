# Spew 打印机

[`Spew`](https://github.com/davecgh/go-spew) 一个漂亮的深度打印机

## 快速入门

### 安装

```shell
go get -u github.com/davecgh/go-spew/spew
```

### 使用

```go
// 引入依赖
import "github.com/davecgh/go-spew/spew"

// 向控制台打印
spew.Dump(myVar1, myVar2, ...)
// 向io.Writer打印
spew.Fdump(someWriter, myVar1, myVar2, ...)
// 生成格式字符串
str := spew.Sdump(myVar1, myVar2, ...)
```

## 说明

`spew`函数按照打印方向可分为三大类

- 控制台打印机

  - `spew.Dump`使用最简单的打印机,类似与原生`fmt`包 但除了 Dump 外都会返回打印数量和错误
  - `spew.Print`
  - `spew.Printf`
  - `spew.Println`

- io 设备打印机 (以 F 开头) 返回写入数量和报错信息

  - `spew.Fdump`
  - `spew.Fprint`
  - `spew.Fprintf`
  - `spew.Fprintln`

- 字符串打印机 返回字符串
  - `spew.Sdump`
  - `spew.Sprint`
  - `spew.Sprintf`
  - `spew.Sprintln`

## 格式化

以 f 结尾的格式化打印机
可以用他们将数据结构更紧凑的打印。格式化字符支持` %v``%#v `（带类型）`%+v`（带指针地址）`%+#v`（类型&指针地址）

## 自定义

`spew`同时开放了自定义接口，可以通过自定义`*ConfigState`使用个性化的打印机。目前提供了如下配置。

```text
* Indent
	String to use for each indentation level for Dump functions.
	It is a single space by default.  A popular alternative is "\t".

* MaxDepth
	Maximum number of levels to descend into nested data structures.
	There is no limit by default.

* DisableMethods
	Disables invocation of error and Stringer interface methods.
	Method invocation is enabled by default.

* DisablePointerMethods
	Disables invocation of error and Stringer interface methods on types
	which only accept pointer receivers from non-pointer variables.  This option
	relies on access to the unsafe package, so it will not have any effect when
	running in environments without access to the unsafe package such as Google
	App Engine or with the "safe" build tag specified.
	Pointer method invocation is enabled by default.

* DisablePointerAddresses
	DisablePointerAddresses specifies whether to disable the printing of
	pointer addresses. This is useful when diffing data structures in tests.

* DisableCapacities
	DisableCapacities specifies whether to disable the printing of capacities
	for arrays, slices, maps and channels. This is useful when diffing data
	structures in tests.

* ContinueOnMethod
	Enables recursion into types after invoking error and Stringer interface
	methods. Recursion after method invocation is disabled by default.

* SortKeys
	Specifies map keys should be sorted before being printed. Use
	this to have a more deterministic, diffable output.  Note that
	only native types (bool, int, uint, floats, uintptr and string)
	and types which implement error or Stringer interfaces are supported,
	with other types sorted according to the reflect.Value.String() output
	which guarantees display stability.  Natural map order is used by
	default.

* SpewKeys
	SpewKeys specifies that, as a last resort attempt, map keys should be
	spewed to strings and sorted by those strings.  This is only considered
	if SortKeys is true.
```

> 注意
> spew 中有引用到`unsafe`包来开发其高级功能，使用是应引起注意。

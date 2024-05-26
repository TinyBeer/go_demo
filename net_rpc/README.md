`net/rpc` 包提供一种将对象的可导出方法暴露为 RPC 接口的方法。

# 使用说明

## 条件

导出的方法需要满足以下条件，否则会被忽视掉：

- 方法类型可导出
- 方法可导出
- 方法包含两个参数，包括导出参数和内建参数
- 方法的第二个参数是指针类型
- 返回值中有 error 类型

方法签名形如`func (t *T) MethodName(argType T1, replyType *T2) error`
第一个参数为调用方需要传递的参数，第二个则为需要返回给调用方的结果。第三个参数如果返回非空，而是一个字符串，客户端将认为接收到一个错误。返回错误时，不会向客户端传递`reply`数据。

## 序列化

rpc 默认使用`encoding/gob`进行序列号，如果有需要可以定制`codec`。

## 通信

服务端可以通过一个长连接工作`ServeConn`，但跟常见的是通过服务端持续监听提供服务。客户端通过像 http 一样的方式连接服务端。服务端完成调用任务通过调用结构体的 Done 通道返回结果。

# 示例

## http

- 业务逻辑

```go
type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

```

- 服务端

```go
func server() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
}
```

- 客户端

```go
func client() {
	client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// 同步调用
	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	// 异步调用
	quotient := new(Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // 等带结果 will be equal to divCall

	if replyCall.Error != nil {
		fmt.Println("async call err:", replyCall.Error)
		return
	}
	args = replyCall.Args.(*Args)
	fmt.Printf("Arith.Divide: %d/%d=%d|%d\n", args.A, args.B, quotient.Quo, quotient.Rem)
}
```

- main

```go
package main

import "time"

func main() {
	go server()

	time.Sleep(time.Second)

	client()
}
```

- 运行

```sh
Arith: 7*8=56
Arith.Divide: 7/8=0|7
```

## tcp

简单修改一下就使用 tcp 支持 rpc

- 服务端

```go
func server() {
	arith := new(Arith)
	rpc.Register(arith)
	// rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	// http.Serve(l, nil)
    rpc.Accept(l)
}
```

- 客户端

```go
// client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
client, err := rpc.Dial("tcp", serverAddress+":1234")
if err != nil {
    log.Fatal("dialing:", err)
}
```

# 自定义服务名

默认服务名为注册时使用的对象类型名,使用`RegisterName(name string, rcvr interface{})`注册服务可以自定义

```go
// rpc.Register(arith)
rpc.RegisterName("calc", arith)
```

# 自定义序列化方式

注册的`codec`需要实现`ServerCodec`接口

```go
type ServerCodec interface {
	ReadRequestHeader(*Request) error
	ReadRequestBody(any) error
	WriteResponse(*Response, any) error

	// Close can be called multiple times and must be idempotent.
	Close() error
}
```

- 服务端

注意这里的 ServeCodec 不会创建服务，仅是执行数据流程。也就是连接需要在 your_codec 中维护。

```go
rpc.ServeCodec(youer_codec)
```

同服务端一样需要自行维护连接。总之很不好用。

- 客户端

```go
client := rpc.NewClientWithCodec(your_codec)
```

# 处理一次数据交互

需要用户提供连接，rpc 来执行一次 rpc 调用。

```go
// 调用之前需要服务端注册
rpc.ServeConn(conn)
```

# 参考

[rpc](https://pkg.go.dev/net/rpc)

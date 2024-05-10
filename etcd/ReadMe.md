# ECTD简介
[ETCD](https://etcd.io/)是一个分布式、可靠的键值存储，用于存储分布式系统中最关键的数据。常用作配置中心和注册中心(类似项目有[zookeeper](https://zookeeper.apache.org/)和[consul](https://www.consul.io/))。
ETCD有以下特点:
* 完全复制：集群中的每个节点都可以使用完整的存档
* 高可用性：Etcd可用于避免硬件的单点故障或网络问题
* 一致性：每次读取都会返回跨多主机的最新写入
* 简单：包括一个定义良好、面向用户的API（gRPC）
* 安全：实现了带有可选的客户端证书身份验证的自动化TLS
* 快速：每秒10000次写入的基准速度
* 可靠：使用Raft算法实现了强一致、高可用的服务存储目录

# 安装ETCD
使用docker-compose安装单节点etcd集群,yaml文件docker-compose.yml
```shell
docker-compose up -d
```
# 快速开始
下载官方提供的包`go get go.etcd.io/etcd/client/v3`
## 测试连接
```go
package main

import (
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.16856.101:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Printf("something wriong, type:%T, err:%v", err, err)
		return
	}
	log.Println("connect etcd ok")
	defer cli.Close()
}
```
> 注意：etcd v3使用rpc进行远程过程调用。`client v3`使用`grpc-go`连接etcd。在使用完成后必须确保关闭连接，否则可能存在协程泄露风险。

## API超时限制。
```go
ctx, cancel := context.WithTimeout(context.Background(), timeout)
resp, err := cli.Put(ctx, "sample_key", "sample_value")
cancel()
if err != nil {
    // handle error!
}
// use the response
```
## 写入
```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	resp, err := cli.Put(ctx, "some_key", "some_value")
	cancel()

	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrUserEmpty, rpctypes.ErrAuthFailed:
			log.Fatalf("authorizate error: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v %T", err, err)
		}
	}
```

## 读取
```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	resp, err := cli.Get(ctx, "some_key")
	cancel()

	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrUserEmpty, rpctypes.ErrAuthFailed:
			log.Fatalf("authorizate error: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v %T", err, err)
		}
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
```
## watch
```go
ch := cli.Watch(context.Background(), "some_key")

for resp := range ch {
    for _, ev := range resp.Events {
        fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
    }
}
```

## lease租约
```go
resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		log.Fatal(err)
	}

	// 5秒钟之后, /nazha/ 这个key就会被移除
	_, err = cli.Put(context.TODO(), "/nazha/", "dsb", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

    // // the key  will be kept forever
	// ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	// if kaerr != nil {
	// 	log.Fatal(kaerr)
	// }
	// go func() {
	// 	for {
	// 		ka := <-ch
	// 		fmt.Println("ttl:", ka.TTL)
	// 	}
	// }()

    // // 4s 后续租一次
    // go func() {
	// 	time.Sleep(time.Second * 4)
	// 	res, kaerr := cli.KeepAliveOnce(context.Background(), resp.ID)
	// 	if kaerr != nil {
	// 		log.Fatal(kaerr)
	// 	}
	// 	fmt.Println("ttl:", res.TTL)
	// }()


	now := time.Now()
	for {
		time.Sleep(time.Second)
		resp, err := cli.Get(context.Background(), "/nazha/")
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(resp.Kvs, time.Since(now))
		}
	}
```
## 分布式锁
```go

```

``` shell
# 1. 生成CA证书和密钥
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=etcd-ca" -days 10000 -out ca.crt
 
# 2. 为etcd服务器生成SSL证书和密钥
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=etcd-server" -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 10000
 
# 3. 为etcd客户端生成SSL证书和密钥
openssl genrsa -out client.key 2048
openssl req -new -key client.key -subj "/CN=etcd-client" -out client.csr
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 10000
 
# 4. 服务器带上IP SANs
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650  -extfile .\extensions.conf

# todo 待验证
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650 -extensions v3_req -extfile '<(printf "subjectAltName = IP:192.168.56.101")'
```

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Username:    "root",
		Password:    "root1234",
		Endpoints:   []string{"192.168.56.101:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Printf("something wriong when new client, type:%T, err:%v", err, err)
		return
	}
	log.Println("connect etcd ok")
	defer cli.Close()

	// 创建两个单独的会话用来演示锁竞争
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, "/my-lock/")

	s2, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s2.Close()
	m2 := concurrency.NewMutex(s2, "/my-lock/")

	// 会话s1获取锁
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s1")

	m2Locked := make(chan struct{})
	go func() {
		defer close(m2Locked)
		// 等待直到会话s1释放了/my-lock/的锁
		if err := m2.Lock(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		time.Sleep(time.Second * 5)
		if err := m1.Unlock(context.TODO()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("released lock for s1")
	}()

	<-m2Locked
	fmt.Println("acquired lock for s2")
}

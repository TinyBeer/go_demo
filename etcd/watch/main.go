package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
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

	ch := cli.Watch(context.Background(), "some_key")

	for resp := range ch {
		for _, ev := range resp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
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
}

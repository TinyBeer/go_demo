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

	go func() {
		time.Sleep(time.Second * 4)
		res, kaerr := cli.KeepAliveOnce(context.Background(), resp.ID)
		if kaerr != nil {
			log.Fatal(kaerr)
		}
		fmt.Println("ttl:", res.TTL)
	}()

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
}

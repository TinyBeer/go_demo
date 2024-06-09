package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mogon访问地址
const uri = "mongodb://192.168.56.101:27017"

func main() {
	opts := options.Client().ApplyURI(uri)
	opts.SetAuth(options.Credential{
		// AuthSource: "admin",  // 指定认证数据库，默认为admin
		Username: "my-page",
		Password: "123456",
	})
	opts.SetMaxPoolSize(5)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// 检查连接情况 超时时间2s
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}
	fmt.Println("connect succes")
}

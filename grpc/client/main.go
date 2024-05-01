package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"learn_grpc/pb"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	addr = flag.String("addr", "127.0.0.1:12480", "address to connect")
	name = flag.String("name", "world", "name to greet")
)

type ClientTokenAuth struct {
}

// GetRequestMetadata implements credentials.PerRPCCredentials.
func (c *ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	//Todo:可以根据访问的uri提供不同的认证信息
	log.Println("uri:", uri)
	return map[string]string{
		"appId":  "tinybeer",
		"appKey": "123123",
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials.
func (c *ClientTokenAuth) RequireTransportSecurity() bool {
	return true
}

var _ credentials.PerRPCCredentials = new(ClientTokenAuth)

func main() {
	flag.Parse()

	creds, _ := credentials.NewClientTLSFromFile(
		`D:\goproject\src\go_code\go_demo\grpc\key\test.pem`,
		"*.tinybeer.com")

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(new(ClientTokenAuth)),
		grpc.WithUnaryInterceptor(unaryInterceptor),
		grpc.WithStreamInterceptor(streamInterceptor),
	)
	// conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	sayHello(client)

	lotsOfReplies(client)

	lotsOfRequest(client)

	bidiHello(client)
}

func sayHello(client pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "k1", "v1", "k2-bin", "v2")
	r, err := client.SayHello(ctx, &pb.HelloRequest{
		Name: *name,
	})
	if err != nil {
		log.Println("failed to say hello", err)
		return
	}
	log.Println("greeting:", r.GetReply())
}

func lotsOfReplies(client pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.LotsOfReplies(ctx, &pb.HelloRequest{
		Name: *name,
	})
	if err != nil {
		log.Fatalf("c.LotsOfReplies failed, err: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("reply over")
			break
		}
		if err != nil {
			log.Fatalf("c.LotsOfReplies failed, err: %v", err)
		}
		log.Printf("got reply: %q\n", res.GetReply())
	}
}

func lotsOfRequest(client pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.LotsOfRequests(ctx)
	if err != nil {
		log.Fatalf("c.LotsOfGreetings failed, err: %v", err)
	}
	names := []string{"七米", "q1mi", "沙河娜扎"}
	for _, name := range names {
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("c.LotsOfGreetings stream.Send(%v) failed, err: %v", name, err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("c.LotsOfGreetings failed: %v", err)
	}
	log.Printf("got reply: %v", res.GetReply())
}

func bidiHello(client pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	stream, err := client.BidiHello(ctx)
	if err != nil {
		log.Fatalf("c.BidiHello failed, err: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			// 接收服务端返回的响应
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("c.BidiHello stream.Recv() failed, err: %v", err)
			}
			log.Printf("AI：%s\n", in.GetReply())
		}
	}()

	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	for {
		cmd, _ := reader.ReadString('\n') // 读到换行
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		// 将获取到的数据发送至服务端
		if err := stream.Send(&pb.HelloRequest{Name: cmd}); err != nil {
			log.Fatalf("c.BidiHello stream.Send(%v) failed: %v", cmd, err)
		}
	}
	stream.CloseSend()
	<-waitc
}

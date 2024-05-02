package main

/*
	ALTS 是一种用于 Google Cloud Platform (GCP) 上的应用程序之间的安全通信的协议。
	如果您在非GCP环境中尝试使用ALTS，将会收到“ALTS is only supported on GCP”的错误消息，
	这表明ALTS仅在GCP上得到支持。
*/
import (
	"context"
	"learn_grpc/pb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
)

const (
	Addr = "127.0.0.1:8989"
)

func main() {
	go serverStart()

	for {
		time.Sleep(time.Second * 3)
		clientRequest()
	}

}

func clientRequest() {
	opt := alts.DefaultClientOptions()
	altsTc := alts.NewClientCreds(opt)
	conn, _ := grpc.Dial(Addr, grpc.WithTransportCredentials(altsTc))

	cc := pb.NewGreeterClient(conn)

	resp, err := cc.SayHello(context.Background(), &pb.HelloRequest{
		Name: "tiny",
	})
	log.Println(resp, err)
}

func serverStart() {
	opt := alts.DefaultServerOptions()
	altsTc := alts.NewServerCreds(opt)
	s := grpc.NewServer(grpc.Creds(altsTc))

	lis, _ := net.Listen("tcp", ":8989")
	pb.RegisterGreeterServer(s, new(server))

	s.Serve(lis)
}

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println(req)
	return &pb.HelloResponse{
		Reply: "hello " + req.GetName(),
	}, nil
}

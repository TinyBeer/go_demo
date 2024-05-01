package main

import (
	"context"
	"io"
	"learn_grpc/pb"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("no metadata found")
	} else {
		log.Println("metadata", md)
	}

	log.Println("client greeting", req.GetName())
	return &pb.HelloResponse{
		Reply: "Hello," + req.GetName(),
	}, nil
}

func (s *server) LotsOfReplies(req *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	log.Println("accept lots of replies request:", req)
	words := []string{
		"你好",
		"hello",
		"こんにちは",
		"안녕하세요",
	}
	for _, word := range words {
		data := &pb.HelloResponse{
			Reply: word + req.GetName(),
		}
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) LotsOfRequests(stream pb.Greeter_LotsOfRequestsServer) error {
	reply := "你好"
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("greeting over")
			return stream.SendAndClose(&pb.HelloResponse{Reply: reply})
		}
		if err != nil {
			return err
		}

		log.Println("receive greeting:", res.GetName())
		reply += res.GetName()
	}
}

func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Println("greeting over")
			return nil
		}
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println("receive greet:", in.GetName())
		reply := magic(in.GetName())
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

func main() {
	creds, _ := credentials.NewServerTLSFromFile(`D:\goproject\src\go_code\go_demo\grpc\key\test.pem`,
		`D:\goproject\src\go_code\go_demo\grpc\key\test.key`)
	lis, err := net.Listen("tcp", ":12480")
	if err != nil {
		log.Fatal("failed to listen port", err)
	}
	log.Println("listen port[12480] ...")
	defer lis.Close()

	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)
	pb.RegisterGreeterServer(s, new(server))
	log.Println("server start ...")
	err = s.Serve(lis)
	if err != nil {
		log.Fatal("failed to serve", err)
	}
}

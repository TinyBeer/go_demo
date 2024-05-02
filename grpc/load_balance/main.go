package main

import (
	"context"
	"learn_grpc/pb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
)

const (
	Schema   = "tiny"
	Endpoint = "resolver.tiny.com"
)

var addrs = []string{"127.0.0.1:8972", "127.0.0.1:8973", "127.0.0.1:8974"}

type MyResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *MyResolver) ResolveNow(o resolver.ResolveNowOptions) {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*MyResolver) Close() {}

type MyBuilder struct{}

func (*MyBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &MyResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			Endpoint: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*MyBuilder) Scheme() string { return Schema }

func init() {
	// 注册 ResolverBuilder
	resolver.Register(&MyBuilder{})
}

func main() {
	go serverListen(":8972")
	go serverListen(":8973")

	for {
		time.Sleep(time.Second * 3)
		clientRequest("tiny")
	}
}

func clientRequest(name string) {
	conn, _ := grpc.Dial(
		Schema+":///"+Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)

	resp, err := pb.NewGreeterClient(conn).SayHello(context.Background(), &pb.HelloRequest{
		Name: name,
	})
	log.Println(resp.GetReply(), err)
}

func serverListen(addr string) {
	log.Printf("listen port[%s] ...", addr)
	lis, _ := net.Listen("tcp", addr)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{addr: addr})
	s.Serve(lis)
}

type server struct {
	pb.UnimplementedGreeterServer
	addr string
}

// SayHello implements pb.GreeterServer.
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println(metadata.FromIncomingContext(ctx))
	log.Println(s.addr)

	return &pb.HelloResponse{
		Reply: "hello " + req.GetName(),
	}, nil
}

var _ pb.GreeterServer = new(server)

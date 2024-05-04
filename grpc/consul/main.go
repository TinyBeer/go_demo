package main

import (
	"context"
	"errors"
	"fmt"
	"learn_grpc/pb"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	name string
	ip   string
	port int
	pb.UnimplementedGreeterServer
}

func NewServer(port int) *server {
	return &server{
		name: "hello",
		ip:   "127.0.0.1",
		port: port,
	}
}

func (s *server) serviceId() string {
	return fmt.Sprintf("%s-%s-%d", s.name, s.ip, s.port)
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Reply: fmt.Sprintf("%5d: Hello, %s", s.port, req.GetName()),
	}, nil
}

type consul struct {
	registered []string
	client     *api.Client
}

func NewConsul(addr string) (*consul, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &consul{
		client: c,
	}, nil
}

func (c *consul) registerServer(svc *server) {
	srv := &api.AgentServiceRegistration{
		ID:      svc.serviceId(),
		Name:    svc.name,
		Tags:    []string{"tiny", "hello"},
		Address: svc.ip,
		Port:    svc.port,
	}
	c.registered = append(c.registered, svc.serviceId())
	c.client.Agent().ServiceRegister(srv)
}

func (c *consul) deregisterServer() error {
	var err error
	for _, id := range c.registered {
		err = errors.Join(c.client.Agent().ServiceDeregister(id))
	}
	return err
}
func main() {
	c, err := NewConsul("192.168.56.101:8500")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		svc1 := NewServer(8889)
		c.registerServer(svc1)
		runGrpcServer(svc1)
	}()
	go func() {
		svc1 := NewServer(8888)
		c.registerServer(svc1)
		runGrpcServer(svc1)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go func() {
		for {
			<-quit
			fmt.Println("service down")
			c.deregisterServer()
			os.Exit(0)
		}
	}()

	runGrpcClient()

}

func runGrpcServer(svc *server) {
	lis, _ := net.Listen("tcp", ":"+strconv.Itoa(svc.port))
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, svc)
	log.Println("server start ...")
	err := s.Serve(lis)
	if err != nil {
		log.Fatal("failed to serve", err)
	}
}

func runGrpcClient() {
	conn, _ := grpc.Dial("consul://192.168.56.101:8500/hello?wait=14s",
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	for {

		time.Sleep(time.Second * 5)

		r, err := client.SayHello(context.Background(), &pb.HelloRequest{
			Name: "jerry",
		})
		if err != nil {
			log.Println("failed to say hello", err)
			return
		}
		log.Println("greeting:", r.GetReply())
	}
}

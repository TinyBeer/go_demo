package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// unaryInterceptor 服务端一元拦截器
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// authentication (token verification)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	log.Println(md)
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	log.Println("handle req")
	m, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("RPC failed with error %v\n", err)
	}
	return m, err
}

type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// streamInterceptor 服务端流拦截器
func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// authentication (token verification)
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["authorization"]) {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		fmt.Printf("RPC failed with error %v\n", err)
	}
	return err
}

// valid 校验认证信息.
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// 执行token认证的逻辑
	// 这里是为了演示方便简单判断token是否与"some-secret-token"相等
	return token == "some-secret-token"
}

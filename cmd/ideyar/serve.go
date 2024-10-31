package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func serve() error {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.MiddleWareAuth()))

	server := servicer.NewServer()
	server.RegisterUserAPIServer(grpcServer, server)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

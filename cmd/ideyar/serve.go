package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	livenessImp "github.com/mosaic-2/IdeYar-server/internal/servicers/liveness"
	"github.com/mosaic-2/IdeYar-server/pkg/LivenessService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcPort = ":8888"
var httpPort = ":80"

func serve() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	mux := runtime.NewServeMux()

	livenessServer, err := livenessImp.NewServer()
	if err != nil {
		return fmt.Errorf("faild to init liveness server %v", err)
	}

	LivenessService.RegisterLivenessServer(grpcServer, livenessServer)
	err = LivenessService.RegisterLivenessHandlerFromEndpoint(ctx, mux, grpcPort, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		return fmt.Errorf("faild to register http server %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		return fmt.Errorf("failed to serve HTTP server: %v", err)
	}

	return nil
}

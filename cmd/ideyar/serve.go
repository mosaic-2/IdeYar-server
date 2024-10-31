package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	livenessImp "github.com/mosaic-2/IdeYar-server/internal/servicers/liveness"
	"github.com/mosaic-2/IdeYar-server/pkg/LivenessService"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

func serve() error {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	livenessServer, err := livenessImp.NewServer()
	if err != nil {
		return fmt.Errorf("faild to init liveness server %v", err)
	}
	LivenessService.RegisterLivenessServer(grpcServer, livenessServer)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func run() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	// Register gRPC-Gateway handlers
	err := LivenessService.RegisterLivenessHandlerFromEndpoint(ctx, mux, "localhost"+*grpcPort, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	log.Printf("Starting HTTP server on %s...", *httpPort)
	if err := http.ListenAndServe(*httpPort, mux); err != nil {
		log.Fatalf("Failed to serve HTTP server: %v", err)
	}
}

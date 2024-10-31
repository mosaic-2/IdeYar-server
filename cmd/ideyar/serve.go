package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	livenessImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/liveness"
	livenessService "github.com/mosaic-2/IdeYar-server/pkg/LivenessService"
)

var (
	grpcPort = ":8888"
	httpPort = ":80"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := runGRPCServer(); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

	go func() {
		if err := runHTTPServer(ctx); err != nil {
			log.Fatalf("Failed to run HTTP server: %v", err)
		}
	}()

	<-signalCh
	log.Println("Received stop signal, shutting down...")
	cancel()
}

func runGRPCServer() error {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcPort, err)
	}

	grpcServer := grpc.NewServer()

	// register grpc services here
	livenessServer, err := livenessImpl.NewServer()
	if err != nil {
		return fmt.Errorf("failed to initialize liveness server: %w", err)
	}
	livenessService.RegisterLivenessServer(grpcServer, livenessServer)

	log.Printf("Starting gRPC server on %s", grpcPort)
	return grpcServer.Serve(lis)
}

func runHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	// register http services here
	err := livenessService.RegisterLivenessHandlerFromEndpoint(
		ctx,
		mux,
		"localhost"+grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC gateway endpoint: %w", err)
	}

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	log.Printf("Starting HTTP server on %s", httpPort)
	return httpServer.ListenAndServe()
}

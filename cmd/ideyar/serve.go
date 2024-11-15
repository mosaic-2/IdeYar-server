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

	"github.com/mosaic-2/IdeYar-server/internal/servicers/authservice"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/postservice"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	authservicepb "github.com/mosaic-2/IdeYar-server/pkg/authServicepb"
	postsrvicepb "github.com/mosaic-2/IdeYar-server/pkg/postservicepb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mosaic-2/IdeYar-server/internal/config"
	livenessImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/liveness"
	livenessService "github.com/mosaic-2/IdeYar-server/pkg/LivenessService"
)

var (
	grpcPort = ":8080"
	httpPort = ":8888"
)

func serve() error {
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
	return nil
}

func runGRPCServer() error {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcPort, err)
	}

	secretKey := []byte(config.GetSecretKey())

	grpcServer := grpc.NewServer()

	// register grpc services here
	livenessServer, err := livenessImpl.NewServer()
	if err != nil {
		return fmt.Errorf("failed to initialize liveness server: %w", err)
	}
	livenessService.RegisterLivenessServer(grpcServer, livenessServer)

	// authentication service
	authServer, err := authservice.NewServer(secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize auth server: %w", err)
	}
	authservicepb.RegisterAuthServer(grpcServer, authServer)

	// post service
	postServer, err := postservice.NewServer(secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize post server: %w", err)
	}
	postsrvicepb.RegisterPostServer(grpcServer, postServer)

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

	err = authservicepb.RegisterAuthHandlerFromEndpoint(
		ctx,
		mux,
		"localhost"+grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC gateway endpoint: %w", err)
	}

	err = postsrvicepb.RegisterPostHandlerFromEndpoint(
		ctx,
		mux,
		"localhost"+grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC gateway endpoint: %w", err)
	}

	// Set up CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers
		AllowCredentials: true,
	})

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: c.Handler(util.AuthMiddleware(mux)),
	}

	log.Printf("Starting HTTP server on %s", httpPort)
	return httpServer.ListenAndServe()
}

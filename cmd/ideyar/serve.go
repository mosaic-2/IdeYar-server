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
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mosaic-2/IdeYar-server/internal/config"
	authImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/auth"
	livenessImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/liveness"
	userProfileImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/user-profile"
	livenessService "github.com/mosaic-2/IdeYar-server/pkg/LivenessService"
	"github.com/mosaic-2/IdeYar-server/pkg/UserProfileService"
	"github.com/mosaic-2/IdeYar-server/pkg/authService"
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
	authServer, err := authImpl.NewServer(secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize auth server: %w", err)
	}
	authService.RegisterAuthServer(grpcServer, authServer)

	userProfileServer, err := userProfileImpl.NewServer(secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize user profile server: %w", err)
	}
	UserProfileService.RegisterUserProfileServer(grpcServer, userProfileServer)

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

	err = authService.RegisterAuthHandlerFromEndpoint(
		ctx,
		mux,
		"localhost"+grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC gateway endpoint: %w", err)
	}

	err = UserProfileService.RegisterUserProfileHandlerFromEndpoint(
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
		AllowedOrigins:   []string{"http://localhost:3000", "https://back.ideyar-app.ir"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: c.Handler(mux),
	}

	log.Printf("Starting HTTP server on %s", httpPort)
	return httpServer.ListenAndServe()
}

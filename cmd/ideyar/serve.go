package main

import (
	"context"
	"fmt"
	"github.com/mosaic-2/IdeYar-server/internal/interceptor"
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
	"google.golang.org/grpc/metadata"

	"github.com/mosaic-2/IdeYar-server/internal/config"
	"github.com/mosaic-2/IdeYar-server/pkg/LivenessServicePb"
	"github.com/mosaic-2/IdeYar-server/pkg/UserProfileServicePb"
	"github.com/mosaic-2/IdeYar-server/pkg/authServicePb"
	"github.com/mosaic-2/IdeYar-server/pkg/postServicePb"

	authImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/auth"
	livenessImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/liveness"
	postImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/post"
	userProfileImpl "github.com/mosaic-2/IdeYar-server/internal/servicers/user-profile"
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
	LivenessServicePb.RegisterLivenessServer(grpcServer, livenessServer)

	// authentication service
	authServer, err := authImpl.NewServer(secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize auth server: %w", err)
	}
	authServicePb.RegisterAuthServer(grpcServer, authServer)

	// post service
	postServer, err := postImpl.NewServer(secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize post server: %w", err)
	}
	postServicePb.RegisterPostServer(grpcServer, postServer)

	userProfileServer, err := userProfileImpl.NewServer(secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize user profile server: %w", err)
	}
	UserProfileServicePb.RegisterUserProfileServer(grpcServer, userProfileServer)

	log.Printf("Starting gRPC server on %s", grpcPort)
	return grpcServer.Serve(lis)
}

func runHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(customMetadataAnnotator),
	)

	// register http services here
	err := LivenessServicePb.RegisterLivenessHandlerFromEndpoint(
		ctx,
		mux,
		"localhost"+grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC gateway endpoint: %w", err)
	}

	err = authServicePb.RegisterAuthHandlerFromEndpoint(
		ctx,
		mux,
		"localhost"+grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC gateway endpoint: %w", err)
	}

	err = postServicePb.RegisterPostHandlerFromEndpoint(
		ctx,
		mux,
		"localhost"+grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC gateway endpoint: %w", err)
	}

	err = mux.HandlePath("POST", "/api/post-image", postImpl.HandlePostImage)
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/api/image/{image}", postImpl.HandleImage)
	if err != nil {
		return err
	}

	err = UserProfileServicePb.RegisterUserProfileHandlerFromEndpoint(
		ctx,
		mux,
		"localhost"+grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC gateway endpoint: %w", err)
	}

	err = mux.HandlePath("POST", "/api/user-image", userProfileImpl.HandleUserImage)
	if err != nil {
		return err
	}

	// Set up CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://back.ideyar-app.ir", "https://back.ideyar-app.ir"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: c.Handler(interceptor.AuthMiddleware(mux)),
	}

	log.Printf("Starting HTTP server on %s", httpPort)
	return httpServer.ListenAndServe()
}

func customMetadataAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	userID := r.Header.Get("x-user-id")
	if userID != "" {
		return metadata.Pairs("user-id", userID)
	}
	return nil
}

package authservice

import (
	pb "github.com/mosaic-2/IdeYar-server/pkg/authServicepb"
)

type Server struct {
	pb.UnimplementedAuthServer
	hmacSecret []byte
}

func NewServer(secretKey []byte) (*Server, error) {
	return &Server{
		hmacSecret: secretKey,
	}, nil
}

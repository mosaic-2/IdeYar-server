package postservice

import (
	pb "github.com/mosaic-2/IdeYar-server/pkg/postservicepb"
)

type Server struct {
	pb.UnimplementedPostServer
	hmacSecret []byte
}

func NewServer(secretKey []byte) (*Server, error) {
	return &Server{
		hmacSecret: secretKey,
	}, nil
}

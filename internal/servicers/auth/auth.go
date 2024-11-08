package auth

import (
	_ "github.com/lib/pq"
	"github.com/mosaic-2/IdeYar-server/pkg/authpb"
)

type Server struct {
	authpb.UnimplementedAuthServer
	hmacSecret []byte
}

func NewServer(secretKey []byte) (*Server, error) {
	return &Server{
		hmacSecret: secretKey,
	}, nil
}

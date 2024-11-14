package authImpl

import (
	_ "github.com/lib/pq"
	"github.com/mosaic-2/IdeYar-server/pkg/authService"
)

type Server struct {
	authService.UnimplementedAuthServer
	hmacSecret []byte
}

func NewServer(secretKey []byte) (*Server, error) {
	return &Server{
		hmacSecret: secretKey,
	}, nil
}

package userProfileImpl

import (
	"github.com/mosaic-2/IdeYar-server/pkg/UserProfileServicePb"
)

type Server struct {
	UserProfileServicePb.UnsafeUserProfileServer
	hmacSecret []byte
}

func NewServer(secretKey []byte) (*Server, error) {
	return &Server{
		hmacSecret: secretKey,
	}, nil
}

package userProfileImpl

import (
	"github.com/mosaic-2/IdeYar-server/pkg/UserProfileService"
)

type Server struct {
	UserProfileService.UnsafeUserProfileServer
	hmacSecret []byte
}

func NewServer(secretKey []byte) (*Server, error) {
	return &Server{
		hmacSecret: secretKey,
	}, nil
}

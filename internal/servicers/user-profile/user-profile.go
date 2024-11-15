package userProfileImpl

import (
	"context"
	pb "github.com/mosaic-2/IdeYar-server/pkg/UserProfileService"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ChangeEmail(ctx context.Context, in *pb.ChangeEmailRequest) (*pb.ChangeEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeEmail not implemented")
}
func (s *Server) ChangePassword(context.Context, *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (s *Server) UpdateProfileInfo(context.Context, *pb.UpdateProfileInfoRequest) (*pb.UpdateProfileInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfileInfo not implemented")
}

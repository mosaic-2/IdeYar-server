package livenessImpl

import (
	"context"
	"github.com/mosaic-2/IdeYar-server/pkg/LivenessServicePb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	LivenessServicePb.UnimplementedLivenessServer
}

func NewServer() (LivenessServicePb.LivenessServer, error) {
	return &server{}, nil
}

func (s *server) CheckLiveness(ctx context.Context, _ *emptypb.Empty) (*LivenessServicePb.CheckLivenessResponse, error) {
	return &LivenessServicePb.CheckLivenessResponse{
		IsAlive:   true,
		Message:   "IdeYar is alive",
		Timestamp: timestamppb.Now(),
	}, nil
}
